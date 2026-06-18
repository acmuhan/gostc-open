package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/pkg/epay"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/cache"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type RechargeReq struct {
	Amount  string `binding:"required" json:"amount"`
	PayType string `binding:"required" json:"payType"`
}

type RechargeResp struct {
	PayUrl  string `json:"payUrl"`
	TradeNo string `json:"tradeNo"`
}

func getEpayClient() (*epay.Client, error) {
	var cfg model.SystemConfigPay
	cache.GetSystemConfigPay(&cfg)
	if cfg.Enable != "1" {
		return nil, errors.New("支付功能未启用")
	}
	return &epay.Client{
		ApiUrl:     cfg.ApiUrl,
		Pid:        cfg.Pid,
		Key:        cfg.Key,
		PrivateKey: cfg.PrivateKey,
		PublicKey:  cfg.PublicKey,
		Version:    cfg.ApiVersion,
	}, nil
}

func (s *service) Recharge(claims jwt.Claims, req RechargeReq, clientIP string) (*RechargeResp, error) {
	client, err := getEpayClient()
	if err != nil {
		return nil, err
	}
	var cfgBase model.SystemConfigBase
	cache.GetSystemConfigBase(&cfgBase)
	orderNo := commerce.NewOrderNo("R")
	notifyUrl := cfgBase.BaseUrl + "/api/v1/public/pay/notify"
	returnUrl := cfgBase.BaseUrl + "/normal/wallet"
	payUrl, err := client.GetPayUrl(epay.CreateOrderReq{
		OutTradeNo: orderNo,
		Type:       req.PayType,
		Name:       "积分充值",
		Money:      req.Amount,
		NotifyUrl:  notifyUrl,
		ReturnUrl:  returnUrl,
		ClientIP:   clientIP,
		Param:      claims.Code,
	})
	if err != nil {
		return nil, errors.New("创建支付订单失败")
	}
	db, _, log := repository.Get("")
	gdb := db.UnderlyingDB()
	order := map[string]any{
		"order_no":  orderNo,
		"user_code": claims.Code,
		"biz_type":  model.ORDER_BIZ_RECHARGE,
		"pay_type":  req.PayType,
		"amount":    commerce.ParseDecimal(req.Amount),
		"status":    model.ORDER_STATUS_PENDING,
		"remark":    "积分充值",
	}
	if err := gdb.Table("commerce_order").Create(order).Error; err != nil {
		log.Error("创建充值订单失败", zap.Error(err))
		return nil, errors.New("创建订单失败")
	}
	return &RechargeResp{PayUrl: payUrl, TradeNo: orderNo}, nil
}
