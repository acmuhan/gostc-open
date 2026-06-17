package service

import (
	"server/pkg/bean"
	"server/service/common/commerce"
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
