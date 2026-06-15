package model

import "github.com/shopspring/decimal"

const (
	ORDER_STATUS_PENDING  = 1
	ORDER_STATUS_PAID     = 2
	ORDER_STATUS_CLOSED   = 3
	ORDER_STATUS_REFUNDED = 4

	ORDER_PAY_BALANCE = "balance"
	ORDER_PAY_POINTS  = "points"
	ORDER_PAY_FREE    = "free"
	ORDER_PAY_ADMIN   = "admin"

	ORDER_BIZ_TUNNEL_CREATE  = "tunnel_create"
	ORDER_BIZ_TUNNEL_RENEW   = "tunnel_renew"
	ORDER_BIZ_HOST_CREATE    = "host_create"
	ORDER_BIZ_HOST_RENEW     = "host_renew"
	ORDER_BIZ_FORWARD_CREATE = "forward_create"
	ORDER_BIZ_FORWARD_RENEW  = "forward_renew"
	ORDER_BIZ_PROXY_CREATE   = "proxy_create"
	ORDER_BIZ_PROXY_RENEW    = "proxy_renew"
	ORDER_BIZ_P2P_CREATE     = "p2p_create"
	ORDER_BIZ_P2P_RENEW      = "p2p_renew"
	ORDER_BIZ_CDK_REDEEM     = "cdk_redeem"
	ORDER_BIZ_ADMIN_ADJUST   = "admin_adjust"
)

type CommerceOrder struct {
	Base
	OrderNo  string          `gorm:"column:order_no;size:64;uniqueIndex;comment:订单号" json:"orderNo"`
	UserCode string          `gorm:"column:user_code;index;comment:用户编号" json:"userCode"`
	User     SystemUser      `gorm:"foreignKey:UserCode;references:Code" json:"-"`
	BizType  string          `gorm:"column:biz_type;size:64;index;comment:业务类型" json:"bizType"`
	BizCode  string          `gorm:"column:biz_code;size:100;index;comment:业务编号" json:"bizCode"`
	PayType  string          `gorm:"column:pay_type;size:32;index;comment:支付方式" json:"payType"`
	Amount   decimal.Decimal `gorm:"column:amount;default:0;comment:订单金额" json:"amount"`
	Points   decimal.Decimal `gorm:"column:points;default:0;comment:积分金额" json:"points"`
	Status   int             `gorm:"column:status;size:1;default:1;index;comment:订单状态" json:"status"`
	Snapshot Map             `gorm:"column:snapshot;type:text;comment:业务快照" json:"snapshot"`
	PaidAt   int64           `gorm:"column:paid_at;index;comment:支付时间" json:"paidAt"`
	Remark   string          `gorm:"column:remark;comment:备注" json:"remark"`
}
