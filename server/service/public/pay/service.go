package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/epay"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/common/commerce"
	"time"
)

type service struct{}

var Service *service

func (s *service) Notify(params epay.NotifyParams) error {
	var cfg model.SystemConfigPay
	cache.GetSystemConfigPay(&cfg)
	client := &epay.Client{
		ApiUrl:     cfg.ApiUrl,
		Pid:        cfg.Pid,
		Key:        cfg.Key,
		PrivateKey: cfg.PrivateKey,
		PublicKey:  cfg.PublicKey,
		Version:    cfg.ApiVersion,
	}
	if !client.Verify(params.ToMap(), params.Sign) {
		return errors.New("签名验证失败")
	}
	if params.TradeStatus != "TRADE_SUCCESS" {
		return nil
	}
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		gdb := tx.UnderlyingDB()
		var order model.CommerceOrder
		if err := gdb.Table("commerce_order").Where("order_no = ?", params.OutTradeNo).First(&order).Error; err != nil {
			return errors.New("订单不存在")
		}
		if order.Status != model.ORDER_STATUS_PENDING {
			return nil
		}
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(order.UserCode)).First()
		if user == nil {
			return errors.New("用户不存在")
		}
		if err := commerce.AdjustWallet(tx, user, order.Amount, model.WALLET_BIZ_RECHARGE, order.OrderNo, "积分充值"); err != nil {
			log.Error("充值到账失败", zap.Error(err))
			return errors.New("充值到账失败")
		}
		gdb.Table("commerce_order").Where("order_no = ?", params.OutTradeNo).Updates(map[string]any{
			"status":  model.ORDER_STATUS_PAID,
			"paid_at": time.Now().Unix(),
		})
		return nil
	})
}
