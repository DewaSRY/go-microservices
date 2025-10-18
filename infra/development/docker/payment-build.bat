
set GOOS=linux
set GOARCH=amd64
go build -o build/payment-service ./services/payment-service/cmd/main.go