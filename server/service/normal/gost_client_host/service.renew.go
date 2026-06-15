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

		host, _ := tx.GostClientHost.Where(
			tx.GostClientHost.UserCode.Eq(user.Code),
			tx.GostClientHost.Code.Eq(req.Code),
		).First()
		if host == nil {
			return errors.New("操作失败")
		}

		if host.ChargingType == model.GOST_CONFIG_CHARGING_ONLY_ONCE || host.ChargingType == model.GOST_CONFIG_CHARGING_FREE {
			return nil
		}
		expAt := commerce.RenewExpAt(host.ExpAt, host.Cycle)
		if _, err := commerce.PayPackage(tx, user, host.Amount, model.ORDER_BIZ_HOST_RENEW, host.Code, model.Map{"hostCode": host.Code, "cycle": host.Cycle, "amount": host.Amount}, "续费套餐"); err != nil {
			return err
		}
		host.Status = 1
		host.ExpAt = expAt
		if err := tx.GostClientHost.Save(host); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        host.Code,
			Type:        model.GOST_TUNNEL_TYPE_HOST,
			ClientCode:  host.ClientCode,
			UserCode:    host.UserCode,
			NodeCode:    host.NodeCode,
			ChargingTye: host.ChargingType,
			ExpAt:       host.ExpAt,
			Limiter:     host.Limiter,
		})
		return nil
	})
}
