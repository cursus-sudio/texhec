module frontend

go 1.24.1

require backend v0.0.0

require (
	github.com/ogiusek/ioc v1.0.8
	github.com/ogiusek/null v1.1.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	github.com/optimus-hft/lockset v0.1.0 // indirect
)

replace backend => ../backend
