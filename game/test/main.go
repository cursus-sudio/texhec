package main

import "github.com/ogiusek/null"

type EntityId int
type Email string
type Mobile string
type Nick string
type Title string

// here user can own book

type UserEmail struct {
	Email Email
}

type UserMobile struct {
	Mobile Mobile
}

type User struct {
	Nick Nick
	// more
}

type Book struct {
	Title  Title
	UserId EntityId
}

type UserEmailRepository interface {
	GetByEntityId(id EntityId) null.Nullable[UserEmail]
	GetByEmail(email Email) (null.Nullable[EntityId], null.Nullable[UserEmail])
}

type UserEmailService interface {
	SetEmail(entityId EntityId, newMail Email) error
	AllowEmailChange(allowedByEmail Email) error
	ConfirmEmail(email Email) error
}

type userEmailService struct {
	UserEmailRepository UserEmailRepository `injected:"1"`
	// change mail component
	// allow change mail component
	// events
}

func (service *userEmailService) SetEmail(entityId EntityId, newMail Email) error {
	// some validation if we want to require entity to be user or something
	emailNullable := service.UserEmailRepository.GetByEntityId(entityId)
	// if email, ok := emailNullable.Ok(); ok {
	if _, ok := emailNullable.Ok(); ok {
		// create allow change mail
		// event created "allow email change" and mailer runs
		return nil
	}
	// create confirm email change
	// event should be created and mailer should run
	return nil
}

func (service *userEmailService) AllowEmailChange(allowedByEmail Email) error {
	// get entity by user email from allow email change repository
	// if entity do not exists return error
	// else remove component and replace it with confirm email change component
	// on creation of this component event should be triggered and it should send mail
	return nil
}

func (service *userEmailService) ConfirmEmail(email Email) error {
	// get entity by user email from confirm email repository
	// if entity do not exists return error
	// if entity exists remove error and remove or replace user email component
	return nil
}
