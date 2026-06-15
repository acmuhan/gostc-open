package service

import (
	"errors"
	"github.com/shopspring/decimal"
	"server/model"
	"server/repository"
	"server/repository/query"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type AdjustReq struct {
	UserCode    string          `binding:"required" json:"userCode"`
	AccountType string          `json:"accountType"`
	Amount      decimal.Decimal `binding:"required" json:"amount"`
	Remark      string          `json:"remark"`
}

func (s *service) Adjust(req AdjustReq) error {
	db, _, _ := repository.Get("")
	return db.Transaction(func(tx *query.Query) error {
		u, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(req.UserCode)).First()
		if u == nil {
			return errors.New("用户不存在")
		}
		account := req.AccountType
		if account == "" {
			account = model.WALLET_ACCOUNT_BALANCE
		}
		if account != model.WALLET_ACCOUNT_BALANCE && account != model.WALLET_ACCOUNT_POINTS {
			return errors.New("账户类型错误")
		}
		if req.Amount.IsZero() {
			return errors.New("调整金额不能为0")
		}
		return commerce.AdjustWallet(tx, u, account, req.Amount, model.WALLET_BIZ_ADMIN_ADJUST, "", req.Remark)
	})
}
