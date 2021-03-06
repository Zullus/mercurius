package context

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/jade"
	"github.com/go-macaron/session"
	"github.com/go-macaron/toolbox"
	"gopkg.in/macaron.v1"
)

var ctx *Context

type Context struct {
	*macaron.Context
	render  jade.Render
	Session session.Store
	Flash   *session.Flash
	Cache   cache.Cache
	Toolbox toolbox.Toolbox
}

func (ctx *Context) HasError() bool {
	hasErr, ok := ctx.Data["HasError"]
	if !ok {
		return false
	}
	ctx.Flash.ErrorMsg = ctx.Data["ErrorMsg"].(string)
	ctx.Data["flash"] = ctx.Flash
	return hasErr.(bool)
}

func (ctx *Context) withErr(msg string, userForm interface{}) {
	if userForm != nil {
		AssignForm(userForm, ctx.Data)
	}
	ctx.Flash.ErrorMsg = msg
	ctx.Data["flash"] = ctx.Flash
}

func (ctx *Context) RenderWithErr(msg string, tpl string, userForm interface{}) {
	ctx.withErr(msg, userForm)
	ctx.HTML(http.StatusOK, tpl)
}

func (ctx *Context) NativeRenderWithErr(msg string, tpl string, userForm interface{}) {
	ctx.withErr(msg, userForm)
	ctx.NativeHTML(http.StatusOK, tpl)
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, r jade.Render, session session.Store, flash *session.Flash, cache cache.Cache, toolbox toolbox.Toolbox) {
		ctx = &Context{
			Context: c,
			render:  r,
			Session: session,
			Flash:   flash,
			Cache:   cache,
			Toolbox: toolbox,
		}
		c.Map(ctx)
	}
}

func (ctx *Context) HTML(status int, name string) {
	ctx.render.HTML(status, name, ctx.Data)
}

func (ctx *Context) NativeHTML(status int, name string) {
	ctx.Context.HTML(status, name, ctx.Data)
}

func (ctx *Context) JSONWithoutEscape(status int, obj interface{}) {
	ctx.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(&obj)
	if err != nil {
		log.Print("[JSONWithoutEscape]" + err.Error())
		http.Error(ctx.Resp, "{'errors':'JSON Marshaling Error = "+err.Error()+"'}", 500)
		return
	}
	ctx.Status(status)
	log.Println("[JSONWithoutEscape] Returned object: " + string(ret))
	ctx.Resp.Write(ret)
}

func I18n(key string) string {
	return ctx.Tr(key)

}
