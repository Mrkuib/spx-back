package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Mrkuib/spx-back/internal/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/qiniu/go-cdk-driver/kodoblob"
	"gocloud.dev/blob"
	"mime/multipart"
	"os"
	"time"
)

var (
	ErrNotExist = os.ErrNotExist
)

type Config struct {
	Driver string // database driver. default is `mysql`.
	DSN    string // database data source name
	BlobUS string // blob URL scheme
}

type Asset struct {
	ID        string
	Name      string
	AuthorId  string
	Category  string
	IsPublic  int
	Address   string
	AssetType string
	Status    int
	CTime     time.Time
	UTime     time.Time
}

type CodeFile struct {
	ID       string
	Name     string
	AuthorId string
	Address  string
	Ctime    time.Time
	Utime    time.Time
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
func (p *Project) FileInfo(ctx context.Context, id string) (*CodeFile, error) {
	if id != "" {
		var address string
		query := "SELECT address FROM project WHERE id = ?"
		err := p.db.QueryRow(query, id).Scan(&address)
		if err != nil {
			return nil, err
		}
		cloudFile := &CodeFile{
			ID:      id,
			Address: address,
		}
		return cloudFile, nil
	}
	return nil, ErrNotExist
}

// Asset returns an Asset.
func (p *Project) Asset(ctx context.Context, id string) (*Asset, error) {
	asset, err := common.QuerySelectById[Asset](p.db, id)
	if err != nil {
		return nil, err
	}
	err = p.modifyAddress(&asset.Address)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// AssetList list assets
func (p *Project) AssetList(ctx context.Context, pageIndex string, pageSize string, assetType string) (*common.Pagination[Asset], error) {
	wheres := map[string]interface{}{"asset_type": assetType}
	pagination, err := common.QueryByPage[Asset](p.db, pageIndex, pageSize, wheres)
	for i := range pagination.Data {
		err := p.modifyAddress(&pagination.Data[i].Address)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return pagination, nil
}

// modifyAddress transfers relative path to download url
func (p *Project) modifyAddress(address *string) error {
	var data struct {
		Assets    map[string]string `json:"assets"`
		IndexJson string            `json:"indexJson"`
	}
	if err := json.Unmarshal([]byte(*address), &data); err != nil {
		return err
	}
	for key, value := range data.Assets {
		data.Assets[key] = os.Getenv("QINIU_PATH") + value // TODO: Replace with real URL prefix
	}
	if data.IndexJson != "" {
		data.IndexJson = os.Getenv("QINIU_PATH") + data.IndexJson
	}
	modifiedAddress, err := json.Marshal(data)
	if err != nil {
		return err
	}
	*address = string(modifiedAddress)
	return nil
}

func (p *Project) SaveProject(ctx context.Context, codeFile *CodeFile, file multipart.File, header *multipart.FileHeader) (*CodeFile, error) {
	if codeFile.ID == "" {
		path, err := UploadFile(ctx, p, os.Getenv("PROJECT_PATH"), file, header)
		if err != nil {
			return nil, err
		}
		codeFile.Address = path
		codeFile.ID, err = AddProject(p, codeFile)
		return codeFile, err
	} else {
		address := GetProjectAddress(codeFile.ID, p)
		err := p.bucket.Delete(ctx, address)
		if err != nil {
			return nil, err
		}
		path, err := UploadFile(ctx, p, os.Getenv("PROJECT_PATH"), file, header)
		if err != nil {
			return nil, err
		}
		codeFile.Address = path
		return codeFile, UpdateProject(p, codeFile)
	}

}
