package model

type SystemAuditLog struct {
	Base
	UserCode  string `gorm:"column:user_code;size:100;index;comment:操作人编号" json:"userCode"`
	Module    string `gorm:"column:module;size:64;index;comment:操作模块" json:"module"`
	Action    string `gorm:"column:action;size:64;index;comment:操作动作" json:"action"`
	Path      string `gorm:"column:path;size:255;comment:请求路径" json:"path"`
	ReqBody   string `gorm:"column:req_body;type:text;comment:请求内容" json:"reqBody"`
	RespCode  int    `gorm:"column:resp_code;comment:响应码" json:"respCode"`
	IP        string `gorm:"column:ip;size:64;comment:请求IP" json:"ip"`
	UserAgent string `gorm:"column:user_agent;size:512;comment:浏览器标识" json:"userAgent"`
	Ms        int64  `gorm:"column:ms;comment:耗时毫秒" json:"ms"`
}
