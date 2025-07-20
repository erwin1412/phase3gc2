cd payment-service

go test ./internal/infra -v

cd ../product-service

go test ./internal/infra -v

cd ../transaction-service

go test ./internal/infra -v
