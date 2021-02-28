package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
)

func main() {
	etcdWatchKey := "/migration/active"

	endpoint, close := startETCDServer()
	defer close()

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	defer etcd.Close()

	watchChan := etcd.Watch(context.Background(), etcdWatchKey, clientv3.WithPrefix(), clientv3.WithPrevKV())
	fmt.Println("set WATCH on " + etcdWatchKey)

	go func() {
		fmt.Println("started goroutine for PUT...")
		count := 0
		for {
			count++
			migrationID := "migration-" + strconv.Itoa(count)
			etcd.Put(context.Background(), etcdWatchKey+"/"+migrationID, time.Now().String())
			fmt.Println("populated " + etcdWatchKey + "/" + migrationID + " with a value..")
			time.Sleep(5 * time.Second)
			etcd.Delete(context.Background(), etcdWatchKey+"/"+migrationID)
		}

	}()

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			fmt.Printf("Event received! %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			if event.PrevKv != nil {
				fmt.Printf("Prev KV : %q\n", event.PrevKv.Value)
			}
		}
	}
}

func startETCDServer() (endpoint string, close func()) {
	cfg := embed.NewConfig()
	cfg.Logger = "zap"
	cfg.LogOutputs = []string{"/dev/null"}
	cfg.Dir = filepath.Join(os.TempDir(), fmt.Sprint(time.Now().Nanosecond()))

	srv, _ := embed.StartEtcd(cfg)

	select {
	case <-srv.Server.ReadyNotify():
	case <-time.After(3 * time.Second):
		fmt.Println("Failed to start embed.Etcd for tests")
	}

	return cfg.ACUrls[0].String(), func() {
		os.RemoveAll(cfg.Dir)
		srv.Close()
	}
}
