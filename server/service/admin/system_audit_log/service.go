package service

import (
	"server/pkg/bean"
	"server/repository"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type PageReq struct {
	bean.PageParam
	UserCode string `json:"userCode"`
	Module   string `json:"module"`
	Action   string `json:"action"`
}

func (s *service) Page(req PageReq) (any, int64) {
	db, _, _ := repository.Get("")
	var list []map[string]any
	var total int64
	q := commerce.DB(db).Table("system_audit_log")
	if req.UserCode != "" {
		q = q.Where("user_code = ?", req.UserCode)
	}
	if req.Module != "" {
		q = q.Where("module = ?", req.Module)
	}
	if req.Action != "" {
		q = q.Where("action = ?", req.Action)
	}
	q.Count(&total).Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list)
	return list, total
}
