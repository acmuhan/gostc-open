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
	proxy, _ := db.GostClientProxy.Where(db.GostClientProxy.Code.Eq(req.Code), db.GostClientProxy.UserCode.Eq(claims.Code)).First()
	if proxy == nil {
		return errors.New("操作失败")
	}
	if proxy.ChargingType != model.GOST_CONFIG_CHARGING_CUCLE_DAY {
		return errors.New("该计费方式不支持自动续费")
	}
	if req.AutoRenew != model.AUTO_RENEW_ENABLE && req.AutoRenew != model.AUTO_RENEW_DISABLE {
		return errors.New("参数错误")
	}
	proxy.AutoRenew = req.AutoRenew
	return db.GostClientProxy.Save(proxy)
}
