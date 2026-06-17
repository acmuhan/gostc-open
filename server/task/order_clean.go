package task

import (
	"server/model"
	"server/repository"
	"time"
)

// 关闭超过24小时未支付的订单
func orderClean() {
	cutoff := time.Now().Add(-24 * time.Hour)
	db, _, _ := repository.Get("")
	gdb := db.UnderlyingDB()
	gdb.Table("commerce_order").
		Where("status = ? AND created_at < ?", model.ORDER_STATUS_PENDING, cutoff).
		Update("status", model.ORDER_STATUS_CLOSED)
}
