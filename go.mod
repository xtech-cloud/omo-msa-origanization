module omo.msa.organization

go 1.13

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0

require (
	github.com/golang/protobuf v1.4.1
	github.com/labstack/gommon v0.3.0
	github.com/micro/go-micro/v2 v2.3.0
	github.com/micro/go-plugins/config/source/consul/v2 v2.0.3
	github.com/micro/go-plugins/logger/logrus/v2 v2.3.0
	github.com/micro/go-plugins/registry/consul/v2 v2.0.3
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.0.3
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/tidwall/gjson v1.6.0
	go.mongodb.org/mongo-driver v1.3.4
	google.golang.org/protobuf v1.24.0
	gopkg.in/yaml.v2 v2.2.4
)