package service

import (
	"server/repository"
	"server/service/common/commerce"
)

func (s *service) Page(req PageReq) (any, int64) {
	db, _, _ := repository.Get("")
	var list []map[string]any
	var total int64
	q := commerce.DB(db).Table("commerce_order")
	if req.Status > 0 {
		q = q.Where("status = ?", req.Status)
	}
	if req.UserCode != "" {
		q = q.Where("user_code = ?", req.UserCode)
	}
	if req.BizType != "" {
		q = q.Where("biz_type = ?", req.BizType)
	}
	if req.OrderNo != "" {
		q = q.Where("order_no = ?", req.OrderNo)
	}
	q.Count(&total).Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list)
	return list, total
}
