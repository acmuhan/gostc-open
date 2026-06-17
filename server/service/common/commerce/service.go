package commerce

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"server/model"
	"server/repository/query"
)

type PackageSnapshot struct {
	ConfigCode   string          `json:"configCode"`
	ConfigName   string          `json:"configName"`
	ChargingType int             `json:"chargingType"`
	Cycle        int             `json:"cycle"`
	Amount       decimal.Decimal `json:"amount"`
	Limiter      int             `json:"limiter"`
	RLimiter     int             `json:"rLimiter"`
	CLimiter     int             `json:"cLimiter"`
}

func NewOrderNo(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}

func SnapshotConfig(cfg *model.GostNodeConfig) PackageSnapshot {
	return PackageSnapshot{ConfigCode: cfg.Code, ConfigName: cfg.Name, ChargingType: cfg.ChargingType, Cycle: cfg.Cycle, Amount: cfg.Amount, Limiter: cfg.Limiter, RLimiter: cfg.RLimiter, CLimiter: cfg.CLimiter}
}

func ExpAtByConfig(cfg *model.GostNodeConfig, base time.Time) int64 {
	if cfg.ChargingType == model.GOST_CONFIG_CHARGING_CUCLE_DAY {
		return base.Add(time.Duration(cfg.Cycle) * 24 * time.Hour).Unix()
	}
	return base.Unix()
}

func RenewExpAt(current int64, cycle int) int64 {
	base := time.Unix(current, 0)
	if base.Unix() < time.Now().Unix() {
		base = time.Now()
	}
	return base.Add(time.Duration(cycle) * 24 * time.Hour).Unix()
}

func PayPackage(tx *query.Query, user *model.SystemUser, amount decimal.Decimal, bizType string, bizCode string, snapshot model.Map, remark string) (*model.CommerceOrder, error) {
	db := tx.UnderlyingDB()
	order := &model.CommerceOrder{OrderNo: NewOrderNo("O"), UserCode: user.Code, BizType: bizType, BizCode: bizCode, PayType: model.ORDER_PAY_BALANCE, Amount: amount, Status: model.ORDER_STATUS_PENDING, Snapshot: snapshot, Remark: remark}
	if amount.LessThanOrEqual(decimal.Zero) {
		order.PayType = model.ORDER_PAY_FREE
		order.Status = model.ORDER_STATUS_PAID
		order.PaidAt = time.Now().Unix()
		return order, db.Create(order).Error
	}
	if user.Amount.LessThan(amount) {
		return nil, errors.New("积分不足")
	}
	before := user.Amount
	after := before.Sub(amount)
	res := db.Model(&model.SystemUser{}).Where("code = ? AND version = ?", user.Code, user.Version).Updates(map[string]any{"amount": after, "version": user.Version + 1})
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, errors.New("扣减积分失败")
	}
	order.Status = model.ORDER_STATUS_PAID
	order.PaidAt = time.Now().Unix()
	if err := db.Create(order).Error; err != nil {
		return nil, err
	}
	ledger := &model.WalletLedger{UserCode: user.Code, AccountType: model.WALLET_ACCOUNT_AMOUNT, Direction: model.WALLET_DIRECTION_OUT, BizType: model.WALLET_BIZ_ORDER_PAY, BizCode: order.OrderNo, Amount: amount, BalanceBefore: before, BalanceAfter: after, Remark: remark}
	if err := db.Create(ledger).Error; err != nil {
		return nil, err
	}
	user.Amount = after
	user.Version++
	return order, nil
}

func AdjustWallet(tx *query.Query, user *model.SystemUser, value decimal.Decimal, bizType, bizCode, remark string) error {
	if value.IsZero() {
		return nil
	}
	db := tx.UnderlyingDB()
	direction := model.WALLET_DIRECTION_IN
	if value.IsNegative() {
		direction = model.WALLET_DIRECTION_OUT
	}
	amount := value.Abs()
	before := user.Amount
	after := before.Add(value)
	if after.IsNegative() {
		return errors.New("积分不足")
	}
	res := db.Model(&model.SystemUser{}).Where("code = ? AND version = ?", user.Code, user.Version).Updates(map[string]any{"amount": after, "version": user.Version + 1})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("账户已变更，请重试")
	}
	user.Amount = after
	user.Version++
	return db.Create(&model.WalletLedger{UserCode: user.Code, AccountType: model.WALLET_ACCOUNT_AMOUNT, Direction: direction, BizType: bizType, BizCode: bizCode, Amount: amount, BalanceBefore: before, BalanceAfter: after, Remark: remark}).Error
}

func DB(tx *query.Query) *gorm.DB { return tx.UnderlyingDB() }
