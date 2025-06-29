module core

go 1.24.3

require (
	backend v0.0.0
	frontend v0.0.0
	github.com/ogiusek/ioc/v2 v2.0.8
	github.com/ogiusek/null v1.1.0
	github.com/ogiusek/relay/v2 v2.0.6
	github.com/veandco/go-sdl2 v0.4.40
	shared v0.0.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang-migrate/migrate/v4 v4.18.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	github.com/ogiusek/events v1.0.2 // indirect
	github.com/ogiusek/lockset v1.0.1 // indirect
	github.com/optimus-hft/lockset v0.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	lukechampine.com/uint128 v1.3.0 // indirect
	modernc.org/cc/v3 v3.41.0 // indirect
	modernc.org/ccgo/v3 v3.16.15 // indirect
	modernc.org/libc v1.37.6 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.7.2 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/sqlite v1.28.0 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
)

replace frontend => ../engine/frontend

replace backend => ../engine/backend

replace shared => ../engine/shared
