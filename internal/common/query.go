package common

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type FilterCondition struct {
	Column    string      // 列名
	Operation string      // 操作符，如 "=", "<", "!=" ...
	Value     interface{} // 值
}

type Pagination[T any] struct {
	TotalCount int
	TotalPage  int
	Data       []T
}

// QueryByPage 通用的 分页查询
func QueryByPage[T any](db *sql.DB, pageIndexParam string, pageSizeParam string, filters []FilterCondition) (*Pagination[T], error) {
	pageIndex, err := strconv.Atoi(pageIndexParam)
	if err != nil {
		return nil, err
	}
	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil {
		return nil, err
	}
	tableName := getTableName[T]()
	scan := tScan[T]()
	whereClause, args := buildWhereClause(filters)

	var totalCount int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s%s`, tableName, whereClause)
	argsForCount := append([]interface{}{}, args...)
	err = db.QueryRow(countQuery, argsForCount...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	totalPage := (totalCount + pageSize - 1) / pageSize

	offset := (pageIndex - 1) * pageSize
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

// QueryById 通用的 SELECT 查询，唯一查询条件为id
func QueryById[T any](db *sql.DB, id string) (*T, error) {
	wheres := []FilterCondition{{Column: "id", Operation: "=", Value: id}}
	results, err := QuerySelect[T](db, wheres)
	if len(results) == 0 {
		return nil, err
	}
	return &results[0], nil
}

// QuerySelect 通用的 SELECT 查询，可以自定义查询条件
func QuerySelect[T any](db *sql.DB, filters []FilterCondition) ([]T, error) {
	tableName := getTableName[T]()
	scan := tScan[T]()
	whereClause, args := buildWhereClause(filters)

	query := fmt.Sprintf("SELECT * FROM %s%s", tableName, whereClause)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// tScan 创建并返回一个适用于 任意结构体 的Scan
func tScan[T any]() func(rows *sql.Rows) (T, error) {
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

// buildWhereClause 根据 FilterCondition 构建 WHERE 子句
func buildWhereClause(conditions []FilterCondition) (string, []interface{}) {
	var whereClauses []string
	var args []interface{}

	for _, condition := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s %s ?", condition.Column, condition.Operation))
		args = append(args, condition.Value)
	}

	// select undelete result
	whereClauses = append(whereClauses, "status != ?")
	args = append(args, 0)

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = " WHERE " + strings.Join(whereClauses, " AND ") // TODO 支持OR操作
	}

	return whereClause, args
}

// getTableName 基于反射获取表名
func getTableName[T any]() string {
	return strings.ToLower(reflect.TypeOf((*T)(nil)).Elem().Name())
}
