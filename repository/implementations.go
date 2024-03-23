package repository

import (
	"context"
	"log"
)

// GetTestById returns user's name for example function
func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output QueryOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT full_name FROM users WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		log.Println("error querying err:", err)
		return
	}
	return
}

// SignUp fuction to register user account
func (r *Repository) SignUp(ctx context.Context, input UserInput) (output QueryOutput, err error) {
    err = r.Db.QueryRowContext(ctx, "INSERT INTO users (phone_number, full_name, password_hash) VALUES ($1, $2, $3) RETURNING id", input.PhoneNumber, input.FullName, input.Password).Scan(&output.ID)
    if err != nil {
        log.Println("error querying sign up user err:", err)
        return
    }

	return
}

// GetUserData fuction to get user account information
func (r *Repository) GetUserData(ctx context.Context, input UserInput) (output QueryOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT full_name,password_hash FROM users WHERE phone_number = $1", input.PhoneNumber).Scan(&output.Name, &output.Password)
	if err != nil {
		log.Println("error querying get user data err:", err)
		return
	}
	return
}

// UpdateName fuction to update user name
func (r *Repository) UpdateName(ctx context.Context, newName string, oldName string, phoneNumber string) (err error) {
	_, err = r.Db.ExecContext(ctx, "UPDATE users SET full_name = $1 WHERE full_name = $2 AND phone_number = $3", newName, oldName, phoneNumber)
	if err != nil {
		log.Println("error querying update name err:", err)
		return err
	}

	return
}

// UpdatePhoneNumber function to update user number
func (r *Repository) UpdatePhoneNumber(ctx context.Context, newNumber string, oldNumber string, fullName string) (err error) {
	_, err = r.Db.ExecContext(ctx, "UPDATE users SET phone_number = $1 WHERE phone_number = $2 AND full_name = $3", newNumber, oldNumber, fullName)
	if err != nil {
		log.Println("error querying update phone number err : ", err)
		return
	}
	return
}

// Logged function to increment user loggin count
func (r *Repository) Logged(ctx context.Context, phoneNumber string) error {
    _, err := r.Db.ExecContext(ctx, "UPDATE users SET successful_login = successful_login + 1 WHERE phone_number = $1", phoneNumber)
    if err != nil {
		log.Println("error querying increment login successful err:", err)
        return err
    }
    return nil
}