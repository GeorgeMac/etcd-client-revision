Etcd Consistency Assumption
---------------------------

The purpose of this project is to demonstrate an assumption made regarding consistency when using the etcd API and the client-side load balancing strategy.

@influxdata we depend on etcd and we recently experimented with the in go client-side load balancing strategy.

This caused panics in out code and ultimately an outage. The panic was sourced from the STM code in the etcd library contained the following message:

```
etcdserver: mvcc: required revision is a future revision
```

We experienced the issue outlined here: https://github.com/etcd-io/etcd/issues/11963

This project forgoes the STM code and demonstrates a minimal reproduction just using the client. It is much the same as how the original issue author reproduced this issue.

## My Confusion

I feel like the exercise being performed here is reasonable and _should not error_, based on my understanding of the etcd API.
I would like to very much understand what is wrong here as it is forces us to either:

a) not use the client-side load balancing
b) force the stm to perform linearized reads on every request

The STM code clearly works the assumption that if the first read is linearized, then subsequent reads should be able to predicate their reads using the same revision.
But clearly the revision returned by the API in the linearizable read could not be present in all nodes in the cluster.

## Run

Reruirements:
- Docker
- Docker Compose

```
docker-compose up
```

You should observer:
```
test_1         | {"level":"warn","ts":"2020-12-15T11:57:53.031Z","caller":"v3@v3.5.0-alpha.22/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"endpoint://client-e9f8c76c-5d5b-42c2-a153-fd41a020ab7b/etcd-node-0:2379","attempt":0,"error":"rpc error: code = OutOfRange desc = etcdserver: mvcc: required revision is a future revision"}
test_1         | panic: etcdserver: mvcc: required revision is a future revision
test_1         |
test_1         | goroutine 1 [running]:
test_1         | main.getFoo(0xb6a740, 0xc0001016c0, 0xc00028eee0, 0x2, 0x2, 0x5c)
test_1         | 	/somewhere/main.go:23 +0xc7
test_1         | main.doRequests(0xb6a740, 0xc0001016c0)
test_1         | 	/somewhere/main.go:33 +0xd3
test_1         | main.main()
test_1         | 	/somewhere/main.go:52 +0x125
etcd-client-revision_test_1 exited with code 2
```
