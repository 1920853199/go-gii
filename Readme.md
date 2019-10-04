## Gii beego 自动化代码生成

#### 1.介绍
Gii 是一个为了协助快速开发 beego 项目而创建的项目，通过 Gii 您可以很容易地为你已存在的数据表在你指定的目录创建 Model 以及 Controller 。它基于 beego 为你写好created ,update,put,已及 delete 等操作方法。
> 注意不能完全依靠 Gii 为你生成的东西，你需要检查一下再进行使用。

#### 2.安装

您可以通过如下的方式安装 bee 工具：

```
go get github.com/1920853199/go-gii 
```
#### 3.使用
```
package main

import (
	"github.com/1920853199/go-gii"
)

func main() {

	source := "xxxx:xxxxxxxx@tcp(127.0.0.1)/abc"
	gii.Column(source,"article","")

	//beego.Run()
}
```

参数介绍
1. 第一个参数 source :数据库连接信息
2. 第二个参数 name 数据表名称
3. 第三个参数 controllerPath 是指定 Controller 生成的路径

> 直接执行这个文件后会在 beego 项目的目录的 models下生成对应数据表的 model,在 controller 下的指定路径生成控制器

结果：

![](http://117.50.7.147:8888/static/uploads/2019091715571981.png)

Controller ```article.go```代码

```
package home

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) List() {


	limit, _ := beego.AppConfig.Int64("limit") // 一页的数量
	page, _ := c.GetInt64("page", 1)           // 页数
	offset := (page - 1) * limit               // 偏移量

	o := orm.NewOrm()
	obj := new(models.Article)

	var data []*models.Article
	qs := o.QueryTable(obj)

	// 获取数据
	_, err := qs.OrderBy("-id").Limit(limit).Offset(offset).All(&data)
	if err != nil {
		c.Abort("404")
	}


	/*c.Data["json"]= &data
	c.ServeJSON()
	c.StopRun()*/


	// 统计
	count, err := qs.Count()
	if err != nil {
		c.Abort("404")
	}

	c.Data["Data"] = &data
	c.Data["Count"] = count
	c.Data["Limit"] = limit
	c.Data["Page"] = page
}

func (c *ArticleController) Put() {
	id, err := c.GetInt("id", 0)

	if id == 0 {
		c.Abort("404")
	}

	// 基础数据
	o := orm.NewOrm()
	obj := new(models.Article)
	var data []*models.Article
	qs := o.QueryTable(obj)
	err = qs.Filter("id", id).One(&data)
	if err != nil {
		c.Abort("404")
	}
	c.Data["Data"] = data[0]

}

func (c *ArticleController) Update() {

	id, _ := c.GetInt("id", 0)


	/*c.Data["json"] = c.Input()
	c.ServeJSON()
	c.StopRun()*/

	response := make(map[string]interface{})

	o := orm.NewOrm()

	obj := models.Article{Id: id}
	if o.Read(&obj) == nil {
		// 需要补充修改的信息
		// 如 ：obj.Reply = reply

		valid := validation.Validation{}

		// 补充需要验证的信息
		// 如：valid.Required(message.Reply, "Reply")

		if valid.HasErrors() {
			// 如果有错误信息
