module github.com/micro-kit/micro-common

go 1.14

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	github.com/coreos/etcd/mvcc/mvccpb => go.etcd.io/etcd/mvcc/mvccpb v3.3.20+incompatible
	go.etcd.io/bbolt => github.com/coreos/bbolt v1.3.4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0 // grpc对etcd依赖问题
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.115
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/etcd v3.3.20+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.3 // indirect
	github.com/jinzhu/gorm v1.9.12
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/naoina/toml v0.1.1
	github.com/nsqio/go-nsq v1.0.8
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.5.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v0.0.0-00010101000000-000000000000 // indirect
	go.etcd.io/etcd v3.3.20+incompatible
	go.uber.org/zap v1.14.1
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/grpc v1.28.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/redis.v5 v5.2.9
	sigs.k8s.io/yaml v1.2.0 // indirect
)
