package gii

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type column struct {
	name  	string
	value 	string
}

func Column(soure,table string){

	//db, err := sql.Open("mysql", "root:2014gaokao@tcp(172.16.230.140)/abc")
	db, err := sql.Open("mysql", soure)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query("select * from "+table+" limit 1")

	if err != nil {
		panic(err.Error())
	}
	
	types, _ := result.ColumnTypes()

	columnArr 	:= []column{}
	packet := make(map[string]string)
	for _ ,v:= range types {

		name := v.Name()//字段名称
		typeName := v.DatabaseTypeName()//字段类型
		typeName = strings.ToLower(typeName)

		if typeName == "datetime" {
			if name == "created" || name == "updated" {
				typeName = `time.Time		` + "`orm:\"auto_now_add;type(datetime)\"`"
			}else if name == "Updated" {
				typeName = `time.Time		` + "`orm:\"auto_now;type(datetime)\"`"
			}else {
				typeName = `time.Time		` + "`orm:\"type(datetime)\"`"
			}
			packet["time"] = "time"
		}else if typeName == "varchar" || typeName == "char" {
			typeName = "string"
		}else if typeName == "text" || typeName == "longtext" {
			typeName = `string		` + "`orm:\"type(text)\"`"
		}else if typeName == "decimal" {
			typeName = "float64"
		}else if typeName == "int" {}else {
			typeName = "string"
		}


		split := strings.Split(name, "_")

		var key string
		for _,item := range split{
			key += strings.Title(item)
		}

		columnArr = append(columnArr,column{key,typeName})
		
	}
	modelName := strings.Title(table[0:1])+table[1:]
	options := Options{path: "/models", column: columnArr,name:modelName,packet:packet}

	createModel(options)
	options = Options{name:table}
	createController(options)
}
