module github.com/Mas4trt/microservices/payment

go 1.26.5

require (
	github.com/Mas4trt/microservices/shared v0.0.0-20260811182637-77ccaa050a30
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.83.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Mas4trt/microservices/shared => ../shared
