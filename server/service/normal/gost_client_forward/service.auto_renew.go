package service

import (
	"errors"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
)

type AutoRenewReq struct {
	Code      string `binding:"required" json:"code"`
	AutoRenew int    `binding:"required" json:"autoRenew"`
}

func (service *service) AutoRenew(claims jwt.Claims, req AutoRenewReq) error {
	db, _, _ := repository.Get("")
	forward, _ := db.GostClientForward.Where(db.GostClientForward.Code.Eq(req.Code), db.GostClientForward.UserCode.Eq(claims.Code)).First()
	if forward == nil {
		return errors.New("操作失败")
	}
	if forward.ChargingType != model.GOST_CONFIG_CHARGING_CUCLE_DAY {
		return errors.New("该计费方式不支持自动续费")
	}
	if req.AutoRenew != model.AUTO_RENEW_ENABLE && req.AutoRenew != model.AUTO_RENEW_DISABLE {
		return errors.New("参数错误")
	}
	forward.AutoRenew = req.AutoRenew
	return db.GostClientForward.Save(forward)
}
