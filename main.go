package main

import (
	"context"
	"log"
	"time"

	client "github.com/ptabor/etcd/client/v3"
)

func putFoo(cli client.KV) int64 {
	resp, err := cli.Put(context.Background(), "foo", "bar")
	if err != nil {
		panic(err)
	}

	return resp.Header.Revision
}

func getFoo(cli client.KV, opts ...client.OpOption) int64 {
	resp, err := cli.Get(context.Background(), "foo", opts...)
	if err != nil {
		panic(err)
	}

	return resp.Header.Revision
}

func doRequests(cli client.KV) {
	// get foo (linearizable read)
	rev := getFoo(cli)
	// get foo again at same revision (serializable read)
	getFoo(
		cli,
		client.WithRev(rev),
		client.WithSerializable(),
	)
	// update foo
	putFoo(cli)
}

func main() {
	cli, err := client.New(client.Config{
		Endpoints:   []string{"etcd-node-0:2379", "etcd-node-1:2379", "etcd-node-2:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("create etcd client failed.")
	}

	for {
		doRequests(cli)
	}
}
