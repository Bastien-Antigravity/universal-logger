module universal-logger

go 1.25.4

replace github.com/Bastien-Antigravity/flexible-logger => ../flexible-logger

require (
	github.com/Bastien-Antigravity/distributed-config v1.4.1-0.20260403170406-daab1400a5ea
	github.com/Bastien-Antigravity/flexible-logger v0.0.0-00010101000000-000000000000
)

require (
	capnproto.org/go/capnp/v3 v3.1.0-alpha.2 // indirect
	github.com/Bastien-Antigravity/safe-socket v1.4.1 // indirect
	github.com/colega/zeropool v0.0.0-20230505084239-6fb4a4f75381 // indirect
	github.com/edsrzf/mmap-go v1.2.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
