module userscrudmysql

go 1.25.1

replace github.com/tracewayapp/go-lightning => ../../..

replace github.com/tracewayapp/go-lightning/lmy => ../..

require github.com/tracewayapp/go-lightning/lmy v0.1.1

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/tracewayapp/go-lightning v0.0.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)
