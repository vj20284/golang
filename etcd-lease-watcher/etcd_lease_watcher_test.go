package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"

	"github.com/stretchr/testify/assert"
)

func TestWatchExpiredEventsReceiveExpiry(t *testing.T) {
	key := keyPrefix + "my-test-key"
	value := "some data"

	endpoint, close := startETCDServer(t)
	defer close()

	cfg := clientv3.Config{Endpoints: []string{endpoint}}
	cli, err := clientv3.New(cfg)
	assert.Nil(t, err)
	defer cli.Close()
	go watchEvents(cli)

	go func() {
		count := 0
		for {
			_, err = cli.Put(context.Background(), key, value+" "+strconv.Itoa(count))
			fmt.Println("Written event to key " + key)
			assert.Nil(t, err)
			time.Sleep(time.Duration(10 * time.Second))
			count++
		}
	}()
	for {
		time.Sleep(time.Duration(10 * time.Second))
	}
}

// Based on: https://github.com/etcd-io/etcd/blob/v3.4.3/clientv3/snapshot/v3_snapshot_test.go#L38
// Does not work on WSL because of a bug https://github.com/microsoft/WSL/issues/3162, but work everywhere else
func startETCDServer(t *testing.T) (endpoint string, close func()) {
	cfg := embed.NewConfig()
	cfg.Logger = "zap"
	cfg.LogOutputs = []string{"/dev/null"}
	cfg.Dir = filepath.Join(os.TempDir(), fmt.Sprint(time.Now().Nanosecond()))

	srv, err := embed.StartEtcd(cfg)
	assert.Nil(t, err)

	select {
	case <-srv.Server.ReadyNotify():
	case <-time.After(3 * time.Second):
		t.Fatalf("Failed to start embed.Etcd for tests")
	}

	return cfg.ACUrls[0].String(), func() {
		os.RemoveAll(cfg.Dir)
		srv.Close()
	}
}
