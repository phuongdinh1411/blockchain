module github.com/phuongdinh1411/blockchain

go 1.21.3

require github.com/phuongdinh1411/blockchain/core v0.0.0-unpublished

require (
	github.com/boltdb/bolt v1.3.1 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/phuongdinh1411/blockchain/core v0.0.0-unpublished => ./internal/core
