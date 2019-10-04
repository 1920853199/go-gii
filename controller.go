package gii

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func createController(op Options)  {

	op.name = strings.Title(op.name)
	content := `package `
	if op.namespace == ""{
		content += `controllers`
	}else{
		content += op.namespace
	}

	if op.modelsNamespace == "" {
		op.modelsNamespace = "models"
	}
	controllerName := op.modelsNamespace+`.`+op.name



	content +=`

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type `+op.name+`Controller struct {
	beego.Controller
}

func (c *`+op.name+`Controller) List() {


	limit, _ := beego.AppConfig.Int64("limit") // 一页的数量
	page, _ := c.GetInt64("page", 1)           // 页数
	offset := (page - 1) * limit               // 偏移量

	o := orm.NewOrm()
	obj := new(`+controllerName+`)

	var data []*`+controllerName+`
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

func (c *`+op.name+`Controller) Put() {
	id, err := c.GetInt("id", 0)

	if id == 0 {
		c.Abort("404")
	}

	// 基础数据
	o := orm.NewOrm()
	obj := new(`+controllerName+`)
	var data []*`+controllerName+`
	qs := o.QueryTable(obj)
	err = qs.Filter("id", id).One(&data)
	if err != nil {
		c.Abort("404")
	}
	c.Data["Data"] = data[0]

}

func (c *`+op.name+`Controller) Update() {

	id, _ := c.GetInt("id", 0)


	/*c.Data["json"] = c.Input()
	c.ServeJSON()
	c.StopRun()*/

	response := make(map[string]interface{})

	o := orm.NewOrm()

	obj := `+controllerName+`{Id: id}
	if o.Read(&obj) == nil {
		// 需要补充修改的信息
		// 如 ：obj.Reply = reply

		valid := validation.Validation{}

		// 补充需要验证的信息
		// 如：valid.Required(message.Reply, "Reply")

		if valid.HasErrors() {
			// 如果有错误信息，证明验证没通过
			// 打印错误信息
			for _, err := range valid.Errors {
				//log.Println(err.Key, err.Message)
				response["msg"] = "新增失败！"
				response["code"] = 500
				response["err"] = err.Key + " " + err.Message
				c.Data["json"] = response
				c.ServeJSON()
				c.StopRun()
			}
		}

		if _, err := o.Update(&obj); err == nil {
			response["msg"] = "修改成功！"
			response["code"] = 200
			response["id"] = id
		} else {
			response["msg"] = "修改失败！"
			response["code"] = 500
			response["err"] = err.Error()
		}
	} else {
		response["msg"] = "修改失败！"
		response["code"] = 500
		response["err"] = "ID 不能为空！"
	}

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}

func (c *`+op.name+`Controller) Delete() {
	id, _ := c.GetInt("id", 0)

	response := make(map[string]interface{})

	o := orm.NewOrm()
	obj := `+controllerName+`{Id: id}

	if _, err := o.Delete(&obj); err == nil {
		response["msg"] = "删除成功！"
		response["code"] = 200
	}else{
		response["msg"] = "删除失败！"
		response["code"] = 500
		response["err"] = err.Error()
	}

	c.Data["json"] = response
	c.ServeJSON()
	c.StopRun()
}`


	dir, _ := os.Getwd()
	if op.path == "" {
		op.path = "/controllers"
	}else{
		op.path = "/controllers/" + op.path
	}
	path := dir + op.path + `/` + strings.ToLower(op.name) + ".go"


	outputFile, outputError := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.WriteString(content)
	outputWriter.Flush()

}
