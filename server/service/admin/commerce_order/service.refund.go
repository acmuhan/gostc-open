package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/commerce"
)

type RefundReq struct {
	OrderNo string `binding:"required" json:"orderNo"`
	Remark  string `json:"remark"`
}

func (s *service) Refund(req RefundReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		gdb := tx.UnderlyingDB()
		var order model.CommerceOrder
		if err := gdb.Table("commerce_order").Where("order_no = ?", req.OrderNo).First(&order).Error; err != nil {
			return errors.New("订单不存在")
		}
		if order.Status != model.ORDER_STATUS_PAID {
			return errors.New("该订单状态不支持退款")
		}
		if order.PayType == model.ORDER_PAY_FREE {
			return errors.New("免费订单无需退款")
		}
		user, err := tx.SystemUser.Where(tx.SystemUser.Code.Eq(order.UserCode)).First()
		if err != nil {
			return errors.New("用户不存在")
		}
		remark := "订单退款"
		if req.Remark != "" {
			remark = req.Remark
		}
		if err := commerce.AdjustWallet(tx, user, order.Amount, model.WALLET_BIZ_REFUND, order.OrderNo, remark); err != nil {
			log.Error("退款失败", zap.Error(err))
			return errors.New("退款失败")
		}
		if err := gdb.Table("commerce_order").Where("order_no = ?", req.OrderNo).Updates(map[string]any{"status": model.ORDER_STATUS_REFUNDED, "remark": remark}).Error; err != nil {
			log.Error("更新订单状态失败", zap.Error(err))
			return errors.New("退款失败")
		}
		return nil
	})
}
