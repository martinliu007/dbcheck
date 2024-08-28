package model
//定义创建巡检结果的结构体
type DBCheck struct {
	DBType         string
	IPPort          string
	CheckType        string
	CheckField      string
	CheckValues     string
}