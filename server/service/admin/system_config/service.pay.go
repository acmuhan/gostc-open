package service

import (
	"errors"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
)

type PayReq struct {
	Enable     string `binding:"required" json:"enable"`
	ApiVersion string `binding:"required" json:"apiVersion"`
	ApiUrl     string `json:"apiUrl"`
	Pid        string `json:"pid"`
	Key        string `json:"key"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func (service *service) Pay(req PayReq) error {
	db, _, log := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		_, _ = tx.SystemConfig.Where(tx.SystemConfig.Kind.Eq(model.SYSTEM_CONFIG_KIND_PAY)).Delete()
		if err := tx.SystemConfig.Create(model.GenerateSystemConfigPay(
			req.Enable,
			req.ApiVersion,
			req.ApiUrl,
			req.Pid,
			req.Key,
			req.PrivateKey,
			req.PublicKey,
		)...); err != nil {
			log.Error("修改支付配置失败", zap.Error(err))
			return errors.New("操作失败")
		}
		cache.SetSystemConfigPay(model.SystemConfigPay{
			Enable:     req.Enable,
			ApiVersion: req.ApiVersion,
			ApiUrl:     req.ApiUrl,
			Pid:        req.Pid,
			Key:        req.Key,
			PrivateKey: req.PrivateKey,
			PublicKey:  req.PublicKey,
		})
		return nil
	})
}
