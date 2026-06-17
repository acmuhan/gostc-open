package model

import "github.com/shopspring/decimal"

const (
	WALLET_DIRECTION_IN  = 1
	WALLET_DIRECTION_OUT = 2

	WALLET_ACCOUNT_AMOUNT = "amount"

	WALLET_BIZ_RECHARGE     = "recharge"
	WALLET_BIZ_CONSUME      = "consume"
	WALLET_BIZ_REFUND       = "refund"
	WALLET_BIZ_ADMIN_ADJUST = "admin_adjust"
	WALLET_BIZ_CDK_REDEEM   = "cdk_redeem"
	WALLET_BIZ_ORDER_PAY    = "order_pay"
)

type WalletLedger struct {
	Base
	UserCode      string          `gorm:"column:user_code;index;comment:用户编号" json:"userCode"`
	User          SystemUser      `gorm:"foreignKey:UserCode;references:Code" json:"-"`
	AccountType   string          `gorm:"column:account_type;size:32;index;comment:账户类型" json:"accountType"`
	Direction     int             `gorm:"column:direction;size:1;index;comment:收支方向" json:"direction"`
	BizType       string          `gorm:"column:biz_type;size:64;index;comment:业务类型" json:"bizType"`
	BizCode       string          `gorm:"column:biz_code;size:100;index;comment:业务编号" json:"bizCode"`
	Amount        decimal.Decimal `gorm:"column:amount;default:0;comment:变动金额" json:"amount"`
	BalanceBefore decimal.Decimal `gorm:"column:balance_before;default:0;comment:变动前余额" json:"balanceBefore"`
	BalanceAfter  decimal.Decimal `gorm:"column:balance_after;default:0;comment:变动后余额" json:"balanceAfter"`
	Remark        string          `gorm:"column:remark;comment:备注" json:"remark"`
}
