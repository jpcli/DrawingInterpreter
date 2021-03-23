package main

import (
	"DrawingInterpreter/drawer"
	"DrawingInterpreter/lexer"
	"DrawingInterpreter/node"
	"DrawingInterpreter/parser"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// 网页路由
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/interpret", interpret)
	http.HandleFunc("/pic", getPic)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("./view/index.html")
	w.Write(file)
}

type jsonDataStruct map[string]interface{}

type fetchJson struct {
	Statements string `json:"statements"`
}

func interpret(w http.ResponseWriter, r *http.Request) {
	// 出错，返回错误
	defer func() {
		if err := recover(); err != nil {
			j, _ := json.Marshal(jsonDataStruct{
				"code": 0,
				"msg":  err,
			})
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Write(j)
		}
	}()

	// 获取post的json数据
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(fmt.Sprintf("read body err, %v\n", err))
	}
	var fetchData fetchJson
	if err = json.Unmarshal(body, &fetchData); err != nil {
		panic(fmt.Sprintf("Unmarshal err, %v\n", err))
	}
	str := fetchData.Statements

	// 初始化返回数据
	returnData := make(jsonDataStruct)

	// 词法分析，构建token表数据
	tokens := lexer.Lexer(str)
	var l []jsonDataStruct
	i := 0
	for _, v := range tokens {
		// 转指针为文本
		var funcPtr string
		if v.FuncPtr != nil {
			funcPtr = fmt.Sprintf("%p", v.FuncPtr)
		} else {
			funcPtr = "nil"
		}

		// 便于json序列化
		i++
		l = append(l, jsonDataStruct{
			"num":       i,
			"tokenType": v.TokenType,
			"lexeme":    v.Lexeme,
			"value":     strconv.FormatFloat(v.Value, 'g', -1, 64),
			"funcPtr":   funcPtr,
		})
	}
	returnData["tokenData"] = l

	// 语法分析，构建树数据
	statements := parser.Parse(tokens)
	var children []jsonDataStruct
	var tree = jsonDataStruct{"name": "All Statements"}

	for i, value := range statements {
		var rootStatement []jsonDataStruct
		for k, v := range value {
			switch v.(type) {
			case *node.Node:
				rootStatement = append(rootStatement, jsonDataStruct{
					"name":     k,
					"children": []jsonDataStruct{v.(*node.Node).GetTree()},
				})
			}
		}
		children = append(children, jsonDataStruct{
			"name":     fmt.Sprintf("%d", i+1),
			"children": rootStatement,
		})

	}
	tree["children"] = children
	returnData["treeData"] = tree

	// 语义分析：画图，返回图文件
	picName := drawer.Draw(statements)
	returnData["pic"] = "pic?id=" + picName[:len(picName)-4]

	// 返回数据
	returnData["code"] = 1
	j, _ := json.Marshal(returnData)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(j)

}

func getPic(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	path := fmt.Sprintf("./pic/%s.png", values.Get("id"))
	if exist(path) {
		// 读取图片
		file, _ := ioutil.ReadFile(path)
		// 返回
		w.Write(file)
		// 删除图片
		os.Remove(path)

	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func exist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true

}
