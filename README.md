# Functions API

## Local dev

```cmd
dapr run --app-id functions-api --app-protocol grpc --app-port 50001  go run cmd/app/main.go
dapr run --app-id functions-api --log-as-json --app-protocol grpc --app-port 50001  go run cmd/app/main.go
```