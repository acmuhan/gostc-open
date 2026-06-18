package task

import (
	"server/bootstrap"
	"server/global"
)

func init() {
	bootstrap.TaskFunc = func() {
		_, _ = global.Cron.AddFunc("0 0 * * *", gostObs)
		_, _ = global.Cron.AddFunc("0 0 * * *", tunnelClean)  // 每天，清理旧隧道
		_, _ = global.Cron.AddFunc("0 */6 * * *", orderClean) // 每6小时，关闭过期订单
		_, _ = global.Cron.AddFunc("0 1 * * *", autoRenew)    // 每天凌晨1点，自动续费
	}
}
