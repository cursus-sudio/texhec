package clock

import (
	"time"

	"github.com/ogiusek/ioc"
)

// interface

type DateFormat string

func NewDateFormat(date string) DateFormat { return DateFormat(date) }
func (format DateFormat) String() string   { return string(format) }

func (format DateFormat) Parse(date string) (time.Time, error) {
	return time.Parse(format.String(), date)
}
func (format DateFormat) Format(date time.Time) string { return date.Format(format.String()) }

type Clock interface {
	Now() time.Time
}

// impl

type clock struct{}

func (clock *clock) Now() time.Time {
	return time.Now()
}

// package

type Pkg struct {
	dateFormat DateFormat
}

func Package(
	dateFormat DateFormat,
) Pkg {
	return Pkg{
		dateFormat: dateFormat,
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) Clock { return &clock{} })
	ioc.RegisterSingleton(c, func(c ioc.Dic) DateFormat { return pkg.dateFormat })
}
