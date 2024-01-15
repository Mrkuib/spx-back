package core

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"golang.org/x/crypto/scrypt"
)

// UploadFile 上传一个文件到 bucket
func UploadFile(ctx context.Context,p *Project, blobKey string,file multipart.File,header *multipart.FileHeader) (string,error) {
	// 从 header 中获取原始文件名
    originalFilename := header.Filename

    // 提取文件扩展名
    ext := filepath.Ext(originalFilename)

	//文件名加密
	blobKey=blobKey+Encrypt(time.Now().String(), originalFilename)+ext

    // 创建 blob writer	
    w, err := p.bucket.NewWriter(ctx, blobKey, nil)
    if err != nil {
        return "",err
    }
    defer w.Close()

    // 将文件内容复制到 blob writer
    _, err = io.Copy(w, file)
    if err != nil {
        return "",err
    }

    // 关闭 writer 提交文件
    return blobKey,w.Close()
}

func AddSpirit(p *Project,s *Spirit) error{
    sqlStr := "insert into spirit (name,author_id , category, use_counts, is_public, address, create_time,update_time) values (?, ?, ?, ?, ?, ?,?, ?)"
	_, err := p.db.Exec(sqlStr, s.Name, s.AuthorId, s.Category, s.UseCounts, s.IsPublic, s.Address,time.Now(),time.Now())
    if err != nil {
        println(err.Error())
        return err
    }
	return err
}

func Encrypt(salt, password string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	return fmt.Sprintf("%x", string(dk))
}