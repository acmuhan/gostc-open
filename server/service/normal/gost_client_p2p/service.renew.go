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
		p2p, _ := tx.GostClientP2P.Where(tx.GostClientP2P.Code.Eq(req.Code), tx.GostClientP2P.UserCode.Eq(claims.Code)).First()
		if p2p == nil {
			return errors.New("操作失败")
		}

		if p2p.ChargingType == model.GOST_CONFIG_CHARGING_ONLY_ONCE || p2p.ChargingType == model.GOST_CONFIG_CHARGING_FREE {
			return nil
		}
		expAt := commerce.RenewExpAt(p2p.ExpAt, p2p.Cycle)
		if _, err := commerce.PayPackage(tx, user, p2p.Amount, model.ORDER_BIZ_P2P_RENEW, p2p.Code, model.Map{"tunnelCode": p2p.Code, "cycle": p2p.Cycle, "amount": p2p.Amount}, "续费套餐"); err != nil {
			return err
		}
		p2p.Status = 1
		p2p.ExpAt = expAt
		if err := tx.GostClientP2P.Save(p2p); err != nil {
			log.Error("续费用户端口转发失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetTunnelInfo(cache.TunnelInfo{
			Code:        p2p.Code,
			Type:        model.GOST_TUNNEL_TYPE_P2P,
			ClientCode:  p2p.ClientCode,
			UserCode:    p2p.UserCode,
			NodeCode:    p2p.NodeCode,
			ChargingTye: p2p.ChargingType,
			ExpAt:       p2p.ExpAt,
			Limiter:     p2p.Limiter,
		})
		engine.ClientP2PConfig(tx, p2p.Code)
		return nil
	})
}
