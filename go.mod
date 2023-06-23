module github.com/edison-moreland/SceneEngine

go 1.20

require (
	github.com/d5/tengo/v2 v2.16.0
	github.com/dave/jennifer v1.6.1
	github.com/fsnotify/fsnotify v1.6.0
	github.com/gen2brain/raylib-go/raylib v0.0.0-20230511170620-d84e4cc82f8d
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.uber.org/zap v1.24.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20220908164124-27713097b956 // indirect
)

replace github.com/d5/tengo/v2 => github.com/edison-moreland/tengo/v2 v2.0.0-20230623191204-1af957b8c051
