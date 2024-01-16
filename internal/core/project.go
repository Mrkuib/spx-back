package core

import (
	"context"
	"database/sql"
	"mime/multipart"
	"os"
	"time"
	_ "github.com/qiniu/go-cdk-driver/kodoblob"
	"gocloud.dev/blob"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	ErrNotExist = os.ErrNotExist
)

type Config struct {
	Driver string // database driver. default is `mysql`.
	DSN    string // database data source name
	BlobUS string // blob URL scheme
}

type CloudFile struct {
	ID    string
	Address string
	Ctime time.Time
	Mtime time.Time
}

type Spirit struct{
	ID string
	Name string
	AuthorId string
	Category string
	UseCounts int
	IsPublic int
	Address string
	Ctime time.Time
	Utime time.Time
}
type CodeFile struct{
	ID string
	Name string
	AuthorId string
	Address string
	Ctime time.Time
	Utime time.Time
}


type Project struct {
	bucket *blob.Bucket
	db     *sql.DB
}

func New(ctx context.Context, conf *Config) (ret *Project, err error) {
	_ = godotenv.Load("../.env")
	if conf == nil {
		conf = new(Config)
	}
	driver := conf.Driver
	dsn := conf.DSN
	bus := conf.BlobUS
	if driver == "" {
		driver = "mysql"
	}
	if dsn == "" {
		dsn = os.Getenv("GOP_SPX_DSN")
	}
	if bus == "" {
		bus = os.Getenv("GOP_SPX_BLOBUS")
	}
	bucket, err := blob.OpenBucket(ctx, bus)
	if err != nil {
		println(err.Error())
		return
	}
	
	db, err := sql.Open(driver, dsn)
	if err != nil {
		println(err.Error())
		return
	}
	return &Project{bucket, db}, nil
}


// Find file address from db
func (p *Project) FileInfo(ctx context.Context, id string) (*CloudFile,error) {
	if id != "" {
		var address string
		// table need fill
		//TODO
    	query := "SELECT address FROM table WHERE id = ?"
    	err := p.db.QueryRow(query, id).Scan(&address)
    	if err != nil {
        	return nil, err
    	}
		cloudFile := &CloudFile{
			ID: id,
			Address: address,
		}
		return cloudFile,nil
	}
	return nil, ErrNotExist
}

// Upload spirit file to cloud
func (p *Project) UploadSpirit(ctx context.Context,spirit *Spirit ,file multipart.File,header *multipart.FileHeader) error {
	path,err := UploadFile(ctx, p, os.Getenv("SPIRIT_PATH"), file, header)
	if err!=nil{
		return err
	}
	spirit.Address=path
	err = AddSpirit(p,spirit)
	return err
}

func (p *Project) SaveProject(ctx context.Context,codeFile *CodeFile ,file multipart.File,header *multipart.FileHeader) (*CodeFile,error) {
	if codeFile.ID==""{
		path,err := UploadFile(ctx, p, os.Getenv("PROJECT_PATH"), file, header)
		if err!=nil{
			return nil,err
		}
		codeFile.Address=path
		codeFile.ID,err = AddProject(p,codeFile)
		return codeFile,err
	}else{
		address:=GetProjectAddress(codeFile.ID,p)
		err:=p.bucket.Delete(ctx,address)
		if err!=nil {
			return nil,err
		}
		path,err := UploadFile(ctx, p, os.Getenv("PROJECT_PATH"), file, header)
		if err!=nil{
			return nil,err
		}
		codeFile.Address=path
		return codeFile,UpdateProject(p,codeFile)
	}
	
}
