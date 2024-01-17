package common

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Pagination[T any] struct {
	TotalCount int
	TotalPage  int
	Data       []T
}

// QueryByPage 通用的 分页查询
func QueryByPage[T any](db *sql.DB, pageIndexParam string, pageSizeParam string, filters map[string]interface{}) (*Pagination[T], error) {
	// 转换string --> int
	pageIndex, err := strconv.Atoi(pageIndexParam)
	if err != nil {
		return nil, err
	}
	pageSize, err := strconv.Atoi(pageSizeParam)

	if err != nil {
		return nil, err
	}

	// 根据T获取tableName
	tableName := strings.ToLower(reflect.TypeOf((*T)(nil)).Elem().Name())

	// 创建Scan
	scan := TScan[T]()

	// 构建 WHERE
	var whereClauses []string
	var args []interface{}
	for col, val := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// 计算开始获取记录的位置
	offset := (pageIndex - 1) * pageSize

	// 查询总记录数
	var totalCount int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s%s`, tableName, whereClause)
	argsForCount := append([]interface{}{}, args...) // 复制 args 用于总数查询
	err = db.QueryRow(countQuery, argsForCount...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPage := (totalCount + pageSize - 1) / pageSize

	// 执行分页查询
	query := fmt.Sprintf("SELECT * FROM %s%s LIMIT ?, ?", tableName, whereClause)
	argsForQuery := append(args, offset, pageSize) // 添加 LIMIT 参数
	rows, err := db.Query(query, argsForQuery...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data []T
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	return &Pagination[T]{
		TotalCount: totalCount,
		TotalPage:  totalPage,
		Data:       data,
	}, nil
}

// QuerySelectById 通用的 SELECT 查询，唯一查询条件为id
func QuerySelectById[T any](db *sql.DB, id string) (*T, error) {
	wheres := map[string]interface{}{"id": id}
	results, err := QuerySelect[T](db, wheres)
	if len(results) == 0 {
		return nil, err
	}
	return &results[0], nil
}

// QuerySelect 通用的 SELECT 查询，可以自定义查询条件
func QuerySelect[T any](db *sql.DB, filters map[string]interface{}) ([]T, error) {
	// 根据 T 获取表名
	tableName := strings.ToLower(reflect.TypeOf((*T)(nil)).Elem().Name())

	// 创建 Scan 函数
	scan := TScan[T]()

	// 构建 WHERE 子句
	var whereClauses []string
	var args []interface{}
	for col, val := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}
	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// 执行 SELECT 查询
	query := fmt.Sprintf("SELECT * FROM %s%s", tableName, whereClause)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 读取结果
	var results []T
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

// TScan 创建并返回一个适用于 任意结构体 的Scan
func TScan[T any]() func(rows *sql.Rows) (T, error) {
	return func(rows *sql.Rows) (T, error) {
		var item T
		itemVal := reflect.ValueOf(&item).Elem()
		if itemVal.Kind() != reflect.Struct {
			return item, errors.New("item must be a struct")
		}

		columns, err := rows.Columns()
		if err != nil {
			return item, err
		}

		if len(columns) != itemVal.NumField() {
			return item, errors.New("number of columns does not match struct fields")
		}

		columnPointers := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			field := itemVal.Field(i)
			if !field.CanSet() {
				return item, errors.New("cannot set struct field " + columns[i])
			}
			columnPointers[i] = field.Addr().Interface()
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return item, err
		}

		return item, nil
	}
}
