#!/bin/bash

echo "Starting payment-service..."
(cd payment-service && go run ./cmd/main.go) &

echo "Starting product-service..."
(cd product-service && go run ./cmd/main.go) &

echo "Starting transaction-service..."
(cd transaction-service && go run ./cmd/main.go) &

echo "Starting gateway-service..."
(cd gateway-service && go run ./cmd/main.go) &

wait
