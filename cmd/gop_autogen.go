package main

import (
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
		}{"code": 200, "msg": "OK", "data": map[string]string{"id": res.ID, "address": res.Address}})
	})
//line cmd/project_yap.gox:22:1
	this.Post("/upload/spirit", func(ctx *yap.Context) {
//line cmd/project_yap.gox:23:1
		name := ctx.FormValue("name")
//line cmd/project_yap.gox:24:1
		category := ctx.FormValue("category")
//line cmd/project_yap.gox:25:1
		file, header, _ := ctx.FormFile("file")
//line cmd/project_yap.gox:26:1
		spirit := &core.Spirit{Name: name, AuthorId: "1", Category: category, UseCounts: 0, IsPublic: 0}
//line cmd/project_yap.gox:33:1
		_ = this.p.UploadSpirit(todo, spirit, file, header)
//line cmd/project_yap.gox:34:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "OK"})
	})
//line cmd/project_yap.gox:40:1
	this.Post("/project/save", func(ctx *yap.Context) {
//line cmd/project_yap.gox:41:1
		id := ctx.FormValue("id")
//line cmd/project_yap.gox:42:1
		uid := ctx.FormValue("uid")
//line cmd/project_yap.gox:43:1
		name := ctx.FormValue("name")
//line cmd/project_yap.gox:44:1
		file, header, _ := ctx.FormFile("file")
//line cmd/project_yap.gox:45:1
		codeFile := &core.CodeFile{ID: id, Name: name, AuthorId: uid}
//line cmd/project_yap.gox:50:1
		res, _ := this.p.SaveProject(todo, codeFile, file, header)
//line cmd/project_yap.gox:51:1
		ctx.Json__1(map[string]interface {
		}{"code": 200, "msg": "ok", "data": map[string]string{"id": res.ID, "address": res.Address}})
	})
//line cmd/project_yap.gox:59:1
	conf := &core.Config{}
//line cmd/project_yap.gox:60:1
	this.p, _ = core.New(todo, conf)
//line cmd/project_yap.gox:62:1
	this.Run__1(":8080")
}
func main() {
	yap.Gopt_App_Main(new(project))
}
