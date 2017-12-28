package gmvc

import (
    "sync"
    "html/template"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/frame/ginstance"
)

// 视图对象(一个请求一个视图对象，用完即销毁)
type View struct {
    mu   sync.RWMutex           // 并发互斥锁
    ctl  *Controller            // 所属控制器
    view *gview.View            // 底层视图对象
    data map[string]interface{} // 视图数据
}

// 创建一个MVC请求中使用的视图对象
func NewView(c *Controller) *View {
    return &View{
        ctl  : c,
        view : ginstance.View(),
        data : make(map[string]interface{}),
    }
}

// 批量绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assigns(data map[string]interface{}) {
    view.mu.Lock()
    defer view.mu.Unlock()
    for k, v := range data {
        view.data[k] = v
    }
}

// 绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assign(key string, value interface{}) {
    view.mu.Lock()
    defer view.mu.Unlock()
    view.data[key] = value
}

// 解析模板，并返回解析后的内容
func (view *View) Parse(file string) ([]byte, error) {
    // 查询模板
    tpl, err := view.view.Template(file)
    if err != nil {
        return nil, err
    }
    // 绑定函数
    tpl.BindFunc("include", view.funcInclude)
    // 执行解析
    view.mu.RLock()
    content, err := tpl.Parse(view.data)
    view.mu.RUnlock()
    return content, err
}

// 解析指定模板
func (view *View) Display(files...string) error {
    file := "default"
    if len(files) > 0 {
        file = files[0]
    }
    if content, err := view.Parse(file); err != nil {
        view.ctl.Response.WriteString("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        view.ctl.Response.Write(content)
    }
    return nil
}

// 模板内置方法：include
func (view *View) funcInclude(file string) template.HTML {
    tpl, err := view.view.Template(file)
    if err != nil {
        return template.HTML(err.Error())
    }
    content, err := tpl.Parse(view.data)
    if err != nil {
        return template.HTML(err.Error())
    }
    return template.HTML(content)
}