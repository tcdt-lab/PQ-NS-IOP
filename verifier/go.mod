module verifier

go 1.23.2

replace test.org/cryptography => /home/koosha/Desktop/PQ-NS-IOP/cryptography_helper

replace test.org/protocol => /home/koosha/Desktop/PQ-NS-IOP/protocol

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/stretchr/testify v1.8.1
	go.uber.org/zap v1.27.0
	gopkg.in/yaml.v3 v3.0.1
	test.org/protocol v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.11.1-0.20230711161743-2e82bdd1719d // indirect
	test.org/cryptography v0.0.0-00010101000000-000000000000 // indirect
)
