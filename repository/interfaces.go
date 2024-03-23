// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output QueryOutput, err error)
	SignUp(ctx context.Context, input UserInput) (output QueryOutput, err error)
	GetUserData(ctx context.Context, input UserInput) (output QueryOutput, err error)
	UpdateName(ctx context.Context, newName string, oldName string, phoneNumber string) (err error)
	UpdatePhoneNumber(ctx context.Context, newNumber string, oldNumber string, fullName string) (err error)
	Logged(ctx context.Context, phoneNumber string) (err error)
}
