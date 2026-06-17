package service

import (
	"server/pkg/bean"
)

type service struct{}

var Service *service

type PageReq struct {
	bean.PageParam
	Status   int    `json:"status"`
	UserCode string `json:"userCode"`
	BizType  string `json:"bizType"`
	OrderNo  string `json:"orderNo"`
}
