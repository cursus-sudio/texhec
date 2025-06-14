module github.com/ogiusek/texhec

go 1.24.1

require github.com/ogiusek/ioc v1.0.7

require github.com/ogiusek/null v1.1.0

// require github.com/ogiusek/texhec/frontend v0.0.0

// replace github.com/ogiusek/texhec/frontend => ../frontend

// require github.com/ogiusek/texhec/backend v0.0.0

require (
	backend v0.0.0
	frontend v0.0.0-00010101000000-000000000000
)

replace backend => ../backend

replace frontend => ../frontend
