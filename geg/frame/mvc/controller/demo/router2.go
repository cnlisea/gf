package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/list",        List)
    ghttp.GetServer().BindHandler("/list/page/2", List2)
    ghttp.GetServer().Router.SetRule(`\/list\/page\/(\d+)[\/\?]*`, "/list?page=$1&")
}

func List1(r *ghttp.Request) {
    r.Response.WriteString("list page:" + r.GetQueryString("page"))
}

func List2(r *ghttp.Request) {
    r.Response.WriteString("customed list page")
}