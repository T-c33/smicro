package main

type Option struct {
	Proto3Filename string
	Output         string
	GenClientCode  bool
	GenServerCode  bool
	Prefix         string
	GoPath         string
	ImportFiles    []string
	ProtoPaths     []string
}
