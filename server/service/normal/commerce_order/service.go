package service

import (
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type PageReq struct {
	Page   int `json:"page"`
	Size   int `json:"size"`
	Status int `json:"status"`
}

func (s *service) Page(claims jwt.Claims, req PageReq) (any, int64) {
	db, _, _ := repository.Get("")
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	var list []map[string]any
	var total int64
	q := commerce.DB(db).Table("commerce_order").Where("user_code = ?", claims.Code)
	if req.Status > 0 {
		q = q.Where("status = ?", req.Status)
	}
	q.Count(&total).Order("id desc").Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&list)
	return list, total
}
