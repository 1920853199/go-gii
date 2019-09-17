package gii

import (
	"bufio"
	"fmt"
	"os"
)

type  Options struct {
	name 			string
	path 			string
	column        	[]column
	packet			map[string]string
	namespace   	string
	modelsNamespace		string
}
func createModel(options Options)  {

	content := `package `
	if options.namespace == ""{
		content += "models"
	}else {
		content += options.namespace
	}

	content +=`

import (
	"github.com/astaxie/beego/orm"`

	for _,v := range options.packet {
		content += "\n    "+`"`+v+`"`
	}
	content += "\n" + `)`
	content += "\n\n" +`type ` + options.name + ` struct {`

	for _,vlues := range options.column {

		content += "\n" + `	`+ vlues.name + `			` + vlues.value
	}

	content += "\n}\n"
	content += `
func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(` + options.name + `))
}`

	//fmt.Printf("%s",content)
	dir, _ := os.Getwd()
	if options.path == "" {
		options.path = "/models"
	}
	path := dir + options.path + `/` + options.name + ".go"


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
