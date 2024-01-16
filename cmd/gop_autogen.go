package main

import (
	"context"
	"github.com/Mrkuib/spx-back/internal/core"
	"github.com/goplus/yap"
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
		}{"code": 200, "msg": "OK", "data": map[string]string{"id": res.ID, "address": res.Address}})
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
	this.Post("/project/save", func(ctx *yap.Context) {
//line cmd/project_yap.gox:34:1
		id := ctx.FormValue("id")
//line cmd/project_yap.gox:35:1
		uid := ctx.FormValue("uid")
//line cmd/project_yap.gox:36:1
		name := ctx.FormValue("name")
//line cmd/project_yap.gox:37:1
		file, header, _ := ctx.FormFile("file")
//line cmd/project_yap.gox:38:1
		codeFile := &core.CodeFile{ID: id, Name: name, AuthorId: uid}
//line cmd/project_yap.gox:43:1
		res, _ := this.p.SaveProject(todo, codeFile, file, header)
//line cmd/project_yap.gox:44:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "ok", "data": map[string]string{"id": res.ID, "address": res.Address}})
	})
//line cmd/project_yap.gox:52:1
	conf := &core.Config{}
//line cmd/project_yap.gox:53:1
	this.p, _ = core.New(todo, conf)
//line cmd/project_yap.gox:55:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(project))
}
