package main

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"smicro/util"
)

type MainGenerator struct {
}

func (d *MainGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {

	filename := path.Join(opt.Output, "main/main.go")
	exist := util.IsFileExist(filename)
	if exist {
		return
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file:%s failed, err:%v\n", filename, err)
		return
	}

	defer file.Close()
	err = d.render(file, main_template, metaData)
	if err != nil {
		fmt.Printf("render failed, err:%v\n", err)
		return
	}
	return
}

func (d *MainGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main").Funcs(templateFuncMap)
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	dir := &MainGenerator{}
	RegisterServerGenerator("main generator", dir)
}
