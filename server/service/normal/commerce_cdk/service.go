package service

import (
	"errors"
	"time"

	"server/model"
	"server/pkg/jwt"
	"server/repository"
	"server/repository/query"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type RedeemReq struct {
	Code string `binding:"required" json:"code" label:"兑换码"`
}

func (s *service) Redeem(claims jwt.Claims, req RedeemReq) error {
	db, _, _ := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		gdb := commerce.DB(tx)
		user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(claims.Code)).First()
		if user == nil {
			return errors.New("用户错误")
		}
		var cdk model.CommerceCdk
		if err := gdb.Where("cdk_code = ?", req.Code).First(&cdk).Error; err != nil {
			return errors.New("兑换码不存在")
		}
		if cdk.Status != model.CDK_STATUS_UNUSED {
			return errors.New("兑换码不可用")
		}
		if cdk.ExpiredAt > 0 && cdk.ExpiredAt < time.Now().Unix() {
			return errors.New("兑换码已过期")
		}
		now := time.Now().Unix()
		res := gdb.Model(&model.CommerceCdk{}).Where("cdk_code = ? AND status = ?", req.Code, model.CDK_STATUS_UNUSED).Updates(map[string]any{"status": model.CDK_STATUS_USED, "user_code": user.Code, "used_at": now})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("兑换码不可用")
		}
		if err := commerce.AdjustWallet(tx, user, cdk.Value, model.WALLET_BIZ_CDK_REDEEM, cdk.Code, "CDK兑换"); err != nil {
			return err
		}
		return gdb.Create(&model.CommerceCdkRedeem{CdkCode: cdk.Code, UserCode: user.Code}).Error
	})
}
