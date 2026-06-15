package model

import "github.com/shopspring/decimal"

const (
	CDK_TYPE_BALANCE = "balance"
	CDK_TYPE_POINTS  = "points"

	CDK_STATUS_UNUSED   = 1
	CDK_STATUS_USED     = 2
	CDK_STATUS_DISABLED = 3
)

type CommerceCdk struct {
	Base
	Code      string          `gorm:"column:cdk_code;size:100;uniqueIndex;comment:兑换码" json:"code"`
	Type      string          `gorm:"column:type;size:32;index;comment:兑换类型" json:"type"`
	Value     decimal.Decimal `gorm:"column:value;default:0;comment:兑换数值" json:"value"`
	Status    int             `gorm:"column:status;size:1;default:1;index;comment:状态" json:"status"`
	BatchNo   string          `gorm:"column:batch_no;size:64;index;comment:批次号" json:"batchNo"`
	UserCode  string          `gorm:"column:user_code;index;comment:使用用户" json:"userCode"`
	UsedAt    int64           `gorm:"column:used_at;index;comment:使用时间" json:"usedAt"`
	ExpiredAt int64           `gorm:"column:expired_at;index;comment:过期时间" json:"expiredAt"`
	Remark    string          `gorm:"column:remark;comment:备注" json:"remark"`
}

type CommerceCdkRedeem struct {
	Base
	CdkCode  string      `gorm:"column:cdk_code;size:100;uniqueIndex:uidx_cdk_redeem;comment:兑换码" json:"cdkCode"`
	Cdk      CommerceCdk `gorm:"foreignKey:CdkCode;references:Code" json:"-"`
	UserCode string      `gorm:"column:user_code;uniqueIndex:uidx_cdk_redeem;comment:用户编号" json:"userCode"`
	User     SystemUser  `gorm:"foreignKey:UserCode;references:Code" json:"-"`
}
