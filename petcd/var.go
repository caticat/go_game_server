package petcd

import (
	"context"
	"sync"

	"go.etcd.io/etcd/clientv3"
)

var (
	g_init           sync.Once
	g_client         *clientv3.Client = nil
	g_cfg            *ConfigEtcd      = nil
	g_leaseID        clientv3.LeaseID
	g_leaseCancel    context.CancelFunc = nil
	g_sliCancelWatch []context.CancelFunc
)

func getInit() *sync.Once                       { return &g_init }
func getClient() *clientv3.Client               { return g_client }
func setClient(cli *clientv3.Client)            { g_client = cli }
func getConfig() *ConfigEtcd                    { return g_cfg }
func setConfig(cfg *ConfigEtcd)                 { g_cfg = cfg }
func getLeaseID() clientv3.LeaseID              { return g_leaseID }
func setLeaseID(leaseID clientv3.LeaseID)       { g_leaseID = leaseID }
func getLeaseCancel() context.CancelFunc        { return g_leaseCancel }
func setLeaseCancel(f context.CancelFunc)       { g_leaseCancel = f }
func getSliCancelWatch() []context.CancelFunc   { return g_sliCancelWatch }
func appendSliCancelWatch(c context.CancelFunc) { g_sliCancelWatch = append(g_sliCancelWatch, c) }
