module verifier

go 1.23.2

replace test.org/cryptography => /home/koosha/Desktop/Thesis/impl/PQ-NS-IOP/cryptography_helper

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/stretchr/testify v1.8.1
	go.uber.org/zap v1.27.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)
