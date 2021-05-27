module github.com/arturoguerra/truenas-csi

go 1.16

require (
	github.com/arturoguerra/go-logging v0.0.0-20200217210650-0d44702b7d73
	github.com/container-storage-interface/spec v1.2.0
	github.com/go-playground/validator/v10 v10.6.1
	github.com/golang/protobuf v1.3.1
	github.com/kubernetes-csi/csi-lib-iscsi v0.0.0-20210519140452-fd47a25d3e16
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onrik/logrus v0.9.0 // indirect
	github.com/rexray/gocsi v1.2.2
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/gjson v1.8.0 // indirect
	google.golang.org/grpc v1.19.0
	k8s.io/mount-utils v0.21.1
	k8s.io/utils v0.0.0-20210521133846-da695404a2bc
)
