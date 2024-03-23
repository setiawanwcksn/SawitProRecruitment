// This file contains types that are used in the repository layer.
package repository

type UserInput struct {
	FullName string
	PhoneNumber string
	Password []byte
}

type GetTestByIdInput struct {
	Id string
}

type QueryOutput struct {
	ID int
	Name string
	Password string
	Token string
}
