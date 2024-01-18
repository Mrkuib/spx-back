package main

import (
	"os"
	"github.com/goplus/yap"
	"context"
	"github.com/Mrkuib/spx-back/internal/core"
)

type project struct {
	yap.App
	p *core.Project
}
//line cmd/project_yap.gox:11
func (this *project) MainEntry() {
//line cmd/project_yap.gox:11:1
	todo := context.TODO()
//line cmd/project_yap.gox:13:1
	this.Get("/project/:id", func(ctx *yap.Context) {
//line cmd/project_yap.gox:14:1
		id := ctx.Param("id")
//line cmd/project_yap.gox:15:1
		res, _ := this.p.FileInfo(todo, id)
//line cmd/project_yap.gox:16:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "OK", "data": map[string]string{"id": res.ID, "address": os.Getenv("QINIU_PATH") + res.Address}})
	})
//line cmd/project_yap.gox:23:1
	this.Get("/asset/:id", func(ctx *yap.Context) {
//line cmd/project_yap.gox:24:1
		id := ctx.Param("id")
//line cmd/project_yap.gox:25:1
		asset, _ := this.p.Asset(todo, id)
//line cmd/project_yap.gox:26:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "ok", "data": map[string]*core.Asset{"asset": asset}})
	})
//line cmd/project_yap.gox:33:1
	this.Get("/list/asset/:pageIndex/:pageSize/:assetType", func(ctx *yap.Context) {
//line cmd/project_yap.gox:34:1
		pageIndex := ctx.Param("pageIndex")
//line cmd/project_yap.gox:35:1
		pageSize := ctx.Param("pageSize")
//line cmd/project_yap.gox:36:1
		assetType := ctx.Param("assetType")
//line cmd/project_yap.gox:37:1
		result, _ := this.p.AssetList(todo, pageIndex, pageSize, assetType)
//line cmd/project_yap.gox:38:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "ok", "data": result})
	})
//line cmd/project_yap.gox:45:1
	this.Post("/project/save", func(ctx *yap.Context) {
//line cmd/project_yap.gox:46:1
		id := ctx.FormValue("id")
//line cmd/project_yap.gox:47:1
		uid := ctx.FormValue("uid")
//line cmd/project_yap.gox:48:1
		name := ctx.FormValue("name")
//line cmd/project_yap.gox:49:1
		file, header, _ := ctx.FormFile("file")
//line cmd/project_yap.gox:50:1
		codeFile := &core.CodeFile{ID: id, Name: name, AuthorId: uid}
//line cmd/project_yap.gox:55:1
		res, _ := this.p.SaveProject(todo, codeFile, file, header)
//line cmd/project_yap.gox:56:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "ok", "data": map[string]string{"id": res.ID, "address": os.Getenv("QINIU_PATH") + res.Address}})
	})
//line cmd/project_yap.gox:64:1
	this.Post("/project/fmt", func(ctx *yap.Context) {
//line cmd/project_yap.gox:65:1
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
//line cmd/project_yap.gox:66:1
		ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
//line cmd/project_yap.gox:67:1
		body := ctx.FormValue("body")
//line cmd/project_yap.gox:68:1
		imports := ctx.FormValue("import")
//line cmd/project_yap.gox:69:1
		res := this.p.CodeFmt(todo, body, imports)
//line cmd/project_yap.gox:70:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "ok", "data": res})
	})
//line cmd/project_yap.gox:78:1
	conf := &core.Config{}
//line cmd/project_yap.gox:79:1
	this.p, _ = core.New(todo, conf)
//line cmd/project_yap.gox:81:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(project))
}
