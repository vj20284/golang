package main

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/clientv3"
)

const keyPrefix = "/migration/active/"

func watchEvents(cli *clientv3.Client) <-chan *clientv3.Event {
	events := make(chan *clientv3.Event)

	//wCh := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix(), clientv3.WithPrevKV(), clientv3.WithFilterPut())
	wCh := cli.Watch(context.Background(), keyPrefix)

	go func() {
		defer close(events)
		for wResp := range wCh {
			for _, ev := range wResp.Events {
				fmt.Printf("Event received of type %s with key %q and value %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				events <- ev
			}
		}
	}()

	return events
}
