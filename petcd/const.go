package petcd

import "github.com/coreos/etcd/mvcc/mvccpb"

type ServicePrefix string

const (
	Config  ServicePrefix = "/config"
	Service ServicePrefix = "/service"
)

type funWatchCallback_t func(eventType mvccpb.Event_EventType, prefix string, kv *mvccpb.KeyValue)
