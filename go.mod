module github.com/chaosblade-io/chaosblade-spec-go

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1
	github.com/sirupsen/logrus v1.4.2
	go.uber.org/zap v1.13.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.17.0 // indirect
	sigs.k8s.io/controller-runtime v0.1.12
)

replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
