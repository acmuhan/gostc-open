package service

import (
	"github.com/shopspring/decimal"
	"server/pkg/jwt"
	"server/repository"
	"server/service/common/commerce"
)

type service struct{}

var Service *service

type WalletSummary struct {
	Amount decimal.Decimal `json:"amount"`
}

type LedgerReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (s *service) Summary(claims jwt.Claims) any {
	db, _, _ := repository.Get("")
	user, _ := db.SystemUser.Where(db.SystemUser.Code.Eq(claims.Code)).First()
	if user == nil {
		return WalletSummary{}
	}
	return WalletSummary{
		Amount: user.Amount,
	}
}

func (s *service) Ledger(claims jwt.Claims, req LedgerReq) (any, int64) {
	db, _, _ := repository.Get("")
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}
	var list []map[string]any
	var total int64
	commerce.DB(db).Table("wallet_ledger").Where("user_code = ?", claims.Code).Order("id desc").Count(&total).Offset((req.Page - 1) * req.Size).Limit(req.Size).Find(&list)
	return list, total
}
