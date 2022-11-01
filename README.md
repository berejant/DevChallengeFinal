# DEV Challenge XIX: Backend Final Round - Mine Detectors

Implementation of [Online Final Task Backend | DEV Challenge XIX](https://docs.google.com/document/d/1QuUdFZ3fPTpMuq6sk1urZyVMuxt4nTVOD9-z69jk2M4/edit)


# Stack
1. [Golang](https://go.dev/) - high performance language for make API-layer
2. [Gin Web Framework](https://github.com/gin-gonic/gin) - help build API fast: router, request validation, building response.
3. Built in [Golang image standard library](https://pkg.go.dev/image@go1.19.2)
4. [Postman](https://www.postman.com/) - Useful tool which provide UI for create API tests and next run it within [docker runner](https://hub.docker.com/r/postman/newman/)

## Run app
> docker compose up -d

## See it works:
> curl -i http://127.0.0.1:8080/healthcheck

> curl -i -X POST http://127.0.0.1:8080/api/image-input

## Run tests
> docker compose run postman

Example results:
```
┌─────────────────────────┬───────────────────┬──────────────────┐
│                         │          executed │           failed │
├─────────────────────────┼───────────────────┼──────────────────┤
│              iterations │                 2 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│                requests │                26 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│            test-scripts │                26 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│      prerequest-scripts │                 0 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│              assertions │                96 │                0 │
├─────────────────────────┴───────────────────┴──────────────────┤
│ total run duration: 1070ms                                     │
├────────────────────────────────────────────────────────────────┤
│ total data received: 17.46kB (approx)                          │
├────────────────────────────────────────────────────────────────┤
│ average response time: 10ms [min: 3ms, max: 104ms, s.d.: 18ms] │
└────────────────────────────────────────────────────────────────┘
```

## Run load testing
> docker compose run siege

Note:
 - you can see siege log at `siege/log/siege.log`
 - load testing inside docker on same machine is not so representative. Better to use two hosts: API-server and siege-client.
 - for having complex load testing (closest to real data) we need to prepare image set with high resolutions

Siege result on my machine:
 - Docker resources: 1.4 GHz Quad-Core Intel Core i5, 4 (v)CPUs, RAM 8 GB, 2 GB Swap.
```
Transactions:                  94892 hits
Availability:                 100.00 %
Elapsed time:                  59.97 secs
Data transferred:              10.05 MB
Response time:                  0.01 secs
Transaction rate:            1582.32 trans/sec
Throughput:                     0.17 MB/sec
Concurrency:                    9.08
Successful transactions:       94893
Failed transactions:               0
Longest transaction:            0.56
Shortest transaction:           0.00
```

## Optional. How to view tests.
Import collection `tests/MineDetectors.postman_collection.json` into [Postman app](https://web.postman.co/).
