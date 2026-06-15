package model

import "github.com/shopspring/decimal"

const (
	SYSTEM_IS_ADMIN = 1 // 是管理员
	SYSTEM_NO_ADMIN = 2 // 不是管理员
)

type SystemUser struct {
	Base
	Account       string          `gorm:"column:account;size:100;uniqueIndex;comment:账号"`
	Password      string          `gorm:"column:password;comment:密码"`
	Salt          string          `gorm:"column:salt;size:8;comment:盐"`
	OtpKey        string          `gorm:"column:otp_key;index;default:'';not null"`
	Admin         int             `gorm:"column:admin;size:1;default:2;comment:是否为管理员"`
	Amount        decimal.Decimal `gorm:"column:amount;default:0;comment:积分(兼容旧字段)"`
	Balance       decimal.Decimal `gorm:"column:balance;default:0;comment:余额"`
	Points        decimal.Decimal `gorm:"column:points;default:0;comment:积分"`
	FrozenBalance decimal.Decimal `gorm:"column:frozen_balance;default:0;comment:冻结余额"`
	BindEmail     SystemUserEmail `gorm:"foreignKey:UserCode;references:Code"`
}
