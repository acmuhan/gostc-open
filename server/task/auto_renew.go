package task

import (
	"fmt"
	"go.uber.org/zap"
	"server/model"
	"server/repository"
	"server/repository/cache"
	"server/repository/query"
	"server/service/common/commerce"
	"server/service/engine"
	"time"
)

func autoRenew() {
	db, _, log := repository.Get("")
	now := time.Now().Unix()
	cutoff := time.Now().Add(24 * time.Hour).Unix()
	grace := now - 30*24*3600

	autoRenewTunnels(db, now, cutoff, grace, log)
	autoRenewHosts(db, now, cutoff, grace, log)
	autoRenewForwards(db, now, cutoff, grace, log)
	autoRenewProxys(db, now, cutoff, grace, log)
	autoRenewP2Ps(db, now, cutoff, grace, log)
}

func renewBizCode(tunnelCode string) string {
	return fmt.Sprintf("%s:%d", tunnelCode, time.Now().UnixNano())
}

func autoRenewTunnels(db *query.Query, now, cutoff, grace int64, log *zap.Logger) {
	var tunnels []model.GostClientTunnel
	gdb := db.UnderlyingDB()
	gdb.Table("gost_client_tunnel").Where(
		"charging_type = ? AND auto_renew = ? AND exp_at <= ? AND exp_at >= ?",
		model.GOST_CONFIG_CHARGING_CUCLE_DAY, model.AUTO_RENEW_ENABLE, cutoff, grace,
	).Find(&tunnels)
	for _, tunnel := range tunnels {
		err := db.Transaction(func(tx *query.Query) error {
			user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(tunnel.UserCode)).First()
			if user == nil {
				return fmt.Errorf("用户不存在")
			}
			expAt := commerce.RenewExpAt(tunnel.ExpAt, tunnel.Cycle)
			if _, err := commerce.PayPackage(tx, user, tunnel.Amount, model.ORDER_BIZ_AUTO_RENEW_TUNNEL, renewBizCode(tunnel.Code), model.Map{"tunnelCode": tunnel.Code, "cycle": tunnel.Cycle, "amount": tunnel.Amount}, "自动续费"); err != nil {
				return err
			}
			tdb := tx.UnderlyingDB()
			tdb.Table("gost_client_tunnel").Where("code = ?", tunnel.Code).Updates(map[string]any{"status": 1, "exp_at": expAt})
			cache.SetTunnelInfo(cache.TunnelInfo{
				Code: tunnel.Code, Type: model.GOST_TUNNEL_TYPE_TUNNEL,
				ClientCode: tunnel.ClientCode, UserCode: tunnel.UserCode,
				NodeCode: tunnel.NodeCode, ChargingTye: tunnel.ChargingType,
				ExpAt: expAt, Limiter: tunnel.Limiter,
			})
			engine.ClientTunnelConfig(tx, tunnel.Code)
			return nil
		})
		if err != nil {
			log.Warn("自动续费失败", zap.String("type", "tunnel"), zap.String("code", tunnel.Code), zap.Error(err))
		}
	}
}

func autoRenewHosts(db *query.Query, now, cutoff, grace int64, log *zap.Logger) {
	var hosts []model.GostClientHost
	gdb := db.UnderlyingDB()
	gdb.Table("gost_client_host").Where(
		"charging_type = ? AND auto_renew = ? AND exp_at <= ? AND exp_at >= ?",
		model.GOST_CONFIG_CHARGING_CUCLE_DAY, model.AUTO_RENEW_ENABLE, cutoff, grace,
	).Find(&hosts)
	for _, host := range hosts {
		err := db.Transaction(func(tx *query.Query) error {
			user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(host.UserCode)).First()
			if user == nil {
				return fmt.Errorf("用户不存在")
			}
			expAt := commerce.RenewExpAt(host.ExpAt, host.Cycle)
			if _, err := commerce.PayPackage(tx, user, host.Amount, model.ORDER_BIZ_AUTO_RENEW_HOST, renewBizCode(host.Code), model.Map{"hostCode": host.Code, "cycle": host.Cycle, "amount": host.Amount}, "自动续费"); err != nil {
				return err
			}
			tdb := tx.UnderlyingDB()
			tdb.Table("gost_client_host").Where("code = ?", host.Code).Updates(map[string]any{"status": 1, "exp_at": expAt})
			cache.SetTunnelInfo(cache.TunnelInfo{
				Code: host.Code, Type: model.GOST_TUNNEL_TYPE_HOST,
				ClientCode: host.ClientCode, UserCode: host.UserCode,
				NodeCode: host.NodeCode, ChargingTye: host.ChargingType,
				ExpAt: expAt, Limiter: host.Limiter,
			})
			return nil
		})
		if err != nil {
			log.Warn("自动续费失败", zap.String("type", "host"), zap.String("code", host.Code), zap.Error(err))
		}
	}
}

func autoRenewForwards(db *query.Query, now, cutoff, grace int64, log *zap.Logger) {
	var forwards []model.GostClientForward
	gdb := db.UnderlyingDB()
	gdb.Table("gost_client_forward").Where(
		"charging_type = ? AND auto_renew = ? AND exp_at <= ? AND exp_at >= ?",
		model.GOST_CONFIG_CHARGING_CUCLE_DAY, model.AUTO_RENEW_ENABLE, cutoff, grace,
	).Find(&forwards)
	for _, forward := range forwards {
		err := db.Transaction(func(tx *query.Query) error {
			user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(forward.UserCode)).First()
			if user == nil {
				return fmt.Errorf("用户不存在")
			}
			expAt := commerce.RenewExpAt(forward.ExpAt, forward.Cycle)
			if _, err := commerce.PayPackage(tx, user, forward.Amount, model.ORDER_BIZ_AUTO_RENEW_FORWARD, renewBizCode(forward.Code), model.Map{"forwardCode": forward.Code, "cycle": forward.Cycle, "amount": forward.Amount}, "自动续费"); err != nil {
				return err
			}
			tdb := tx.UnderlyingDB()
			tdb.Table("gost_client_forward").Where("code = ?", forward.Code).Updates(map[string]any{"status": 1, "exp_at": expAt})
			cache.SetTunnelInfo(cache.TunnelInfo{
				Code: forward.Code, Type: model.GOST_TUNNEL_TYPE_FORWARD,
				ClientCode: forward.ClientCode, UserCode: forward.UserCode,
				NodeCode: forward.NodeCode, ChargingTye: forward.ChargingType,
				ExpAt: expAt, Limiter: forward.Limiter,
			})
			engine.ClientForwardConfig(tx, forward.Code)
			return nil
		})
		if err != nil {
			log.Warn("自动续费失败", zap.String("type", "forward"), zap.String("code", forward.Code), zap.Error(err))
		}
	}
}

func autoRenewProxys(db *query.Query, now, cutoff, grace int64, log *zap.Logger) {
	var proxys []model.GostClientProxy
	gdb := db.UnderlyingDB()
	gdb.Table("gost_client_proxy").Where(
		"charging_type = ? AND auto_renew = ? AND exp_at <= ? AND exp_at >= ?",
		model.GOST_CONFIG_CHARGING_CUCLE_DAY, model.AUTO_RENEW_ENABLE, cutoff, grace,
	).Find(&proxys)
	for _, proxy := range proxys {
		err := db.Transaction(func(tx *query.Query) error {
			user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(proxy.UserCode)).First()
			if user == nil {
				return fmt.Errorf("用户不存在")
			}
			expAt := commerce.RenewExpAt(proxy.ExpAt, proxy.Cycle)
			if _, err := commerce.PayPackage(tx, user, proxy.Amount, model.ORDER_BIZ_AUTO_RENEW_PROXY, renewBizCode(proxy.Code), model.Map{"proxyCode": proxy.Code, "cycle": proxy.Cycle, "amount": proxy.Amount}, "自动续费"); err != nil {
				return err
			}
			tdb := tx.UnderlyingDB()
			tdb.Table("gost_client_proxy").Where("code = ?", proxy.Code).Updates(map[string]any{"status": 1, "exp_at": expAt})
			cache.SetTunnelInfo(cache.TunnelInfo{
				Code: proxy.Code, Type: model.GOST_TUNNEL_TYPE_PROXY,
				ClientCode: proxy.ClientCode, UserCode: proxy.UserCode,
				NodeCode: proxy.NodeCode, ChargingTye: proxy.ChargingType,
				ExpAt: expAt, Limiter: proxy.Limiter,
			})
			engine.ClientProxyConfig(tx, proxy.Code)
			return nil
		})
		if err != nil {
			log.Warn("自动续费失败", zap.String("type", "proxy"), zap.String("code", proxy.Code), zap.Error(err))
		}
	}
}

func autoRenewP2Ps(db *query.Query, now, cutoff, grace int64, log *zap.Logger) {
	var p2ps []model.GostClientP2P
	gdb := db.UnderlyingDB()
	gdb.Table("gost_client_p2p").Where(
		"charging_type = ? AND auto_renew = ? AND exp_at <= ? AND exp_at >= ?",
		model.GOST_CONFIG_CHARGING_CUCLE_DAY, model.AUTO_RENEW_ENABLE, cutoff, grace,
	).Find(&p2ps)
	for _, p2p := range p2ps {
		err := db.Transaction(func(tx *query.Query) error {
			user, _ := tx.SystemUser.Where(tx.SystemUser.Code.Eq(p2p.UserCode)).First()
			if user == nil {
				return fmt.Errorf("用户不存在")
			}
			expAt := commerce.RenewExpAt(p2p.ExpAt, p2p.Cycle)
			if _, err := commerce.PayPackage(tx, user, p2p.Amount, model.ORDER_BIZ_AUTO_RENEW_P2P, renewBizCode(p2p.Code), model.Map{"p2pCode": p2p.Code, "cycle": p2p.Cycle, "amount": p2p.Amount}, "自动续费"); err != nil {
				return err
			}
			tdb := tx.UnderlyingDB()
			tdb.Table("gost_client_p2p").Where("code = ?", p2p.Code).Updates(map[string]any{"status": 1, "exp_at": expAt})
			cache.SetTunnelInfo(cache.TunnelInfo{
				Code: p2p.Code, Type: model.GOST_TUNNEL_TYPE_P2P,
				ClientCode: p2p.ClientCode, UserCode: p2p.UserCode,
				NodeCode: p2p.NodeCode, ChargingTye: p2p.ChargingType,
				ExpAt: expAt, Limiter: p2p.Limiter,
			})
			engine.ClientP2PConfig(tx, p2p.Code)
			return nil
		})
		if err != nil {
			log.Warn("自动续费失败", zap.String("type", "p2p"), zap.String("code", p2p.Code), zap.Error(err))
		}
	}
}
