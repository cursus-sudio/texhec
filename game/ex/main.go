package ex

import "github.com/ogiusek/null"

type Id int
type Email string
type Mobile string
type Nick string
type Title string

// here user can own book

type User struct {
	Id     Id
	email  null.Nullable[Email]
	mobile null.Nullable[Mobile]
}

type Book struct {
	Id     Id
	Title  Title
	UserId Id
}
