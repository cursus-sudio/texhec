module core

go 1.24.3

require (
	backend v0.0.0
	frontend v0.0.0
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/mathgl v1.2.0
	github.com/ogiusek/events v1.0.6
	github.com/ogiusek/ioc/v2 v2.0.13
	github.com/ogiusek/relay/v2 v2.0.6
	github.com/veandco/go-sdl2 v0.4.40
	golang.org/x/image v0.32.0
	shared v0.0.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang-migrate/migrate/v4 v4.18.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/optimus-hft/lockset v0.1.0 // indirect
	github.com/optimus-hft/lockset/v2 v2.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	modernc.org/libc v1.66.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.38.2 // indirect
)

replace frontend => ../engine/frontend

replace backend => ../engine/backend

replace shared => ../engine/shared
