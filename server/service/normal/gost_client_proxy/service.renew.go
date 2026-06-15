package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/common/commerce"
	"server/service/engine"
)

type RenewReq struct {
	Code string `binding:"required" json:"code"`
}

func (service *service) Renew(claims jwt.Claims, req RenewReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}

		proxy, _ := tx.GostClientProxy.Where(
			tx.GostClientProxy.UserCode.Eq(user.Code),
			tx.GostClientProxy.Code.Eq(req.Code),
		).First()
		if proxy == nil {
			return errors.New("操作失败")
		}

		if proxy.ChargingType == model.GOST_CONFIG_CHARGING_ONLY_ONCE || proxy.ChargingType == model.GOST_CONFIG_CHARGING_FREE {
			return nil
		}
		expAt := commerce.RenewExpAt(proxy.ExpAt, proxy.Cycle)
		if _, err := commerce.PayPackage(tx, user, proxy.Amount, model.ORDER_BIZ_PROXY_RENEW, proxy.Code, model.Map{"tunnelCode": proxy.Code, "cycle": proxy.Cycle, "amount": proxy.Amount}, "续费套餐"); err != nil {
			return err
		}
		proxy.Status = 1
		proxy.ExpAt = expAt
		if err := tx.GostClientProxy.Save(proxy); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        proxy.Code,
			Type:        model.GOST_TUNNEL_TYPE_PROXY,
			ClientCode:  proxy.ClientCode,
			UserCode:    proxy.UserCode,
			NodeCode:    proxy.NodeCode,
			ChargingTye: proxy.ChargingType,
			ExpAt:       proxy.ExpAt,
			Limiter:     proxy.Limiter,
		})
		engine.ClientProxyConfig(tx, proxy.Code)
		return nil
	})
}
