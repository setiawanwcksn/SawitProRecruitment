package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID          uint
	PhoneNumber string
	Password    string
}

type Claims struct {
	Name        string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.Response
	req := repository.GetTestByIdInput{
		Id: params.Id,
	}
	uid, err := s.Repository.GetTestById(ctx.Request().Context(), req)
	if err != nil {
		resp.Message = "User not found"
		return ctx.JSON(http.StatusBadRequest, resp)
	}
	resp.Message = fmt.Sprintf("Hello User %s", uid.Name)
	return ctx.JSON(http.StatusOK, resp)
}

// GetMyProfile implements generated.ServerInterface.
func (s *Server) GetMyProfile(ctx echo.Context) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" || len("Bearer ") < 1 {
		return echo.NewHTTPError(http.StatusForbidden, "Authorization header not found")
	}
	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("sawit-pro-digital"), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Token is not valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return echo.NewHTTPError(http.StatusForbidden, "Token not found")
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"name":         claims.Name,
		"phone_number": claims.PhoneNumber,
	})
}

// PatchUpdateMyProfile implements generated.ServerInterface.
func (s *Server) UpdateMyProfile(ctx echo.Context, params generated.UpdateMyProfileParams) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusForbidden, "Authorization header tidak ditemukan")
	}
	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("sawit-pro-digital"), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Token tidak valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return echo.NewHTTPError(http.StatusForbidden, "Token tidak valid")
	}

	var newPhoneNumber string
	var newFullName string
	if params.PhoneNumber != nil {
		newPhoneNumber = *params.PhoneNumber
	}
	if params.FullName != nil {
		newFullName = *params.FullName
	}

	if (newPhoneNumber == "" || !validatePhoneNumber(newPhoneNumber)) && (newFullName == "" || !validateFullName(newFullName)) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}

	if newPhoneNumber != "" {
		output, _ := s.Repository.GetUserData(ctx.Request().Context(), repository.UserInput{
			PhoneNumber: newPhoneNumber,
		})

		if output.Name != "" {
			return echo.NewHTTPError(http.StatusConflict, "Phone Number already exists")
		}

		err := s.Repository.UpdatePhoneNumber(ctx.Request().Context(), newPhoneNumber, claims.PhoneNumber, claims.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Update Phone number failed")
		}
	} else if newFullName != "" {
		err := s.Repository.UpdateName(ctx.Request().Context(), newFullName, claims.Name, claims.PhoneNumber)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Update name failed")
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Successfully updated user data",
	})
}

// PostLogin implements generated.ServerInterface.
func (s *Server) PostLogin(ctx echo.Context, params generated.PostLoginParams) error {

	if !validatePhoneNumber(params.PhoneNumber) || !validatePassword(params.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request data")
	}
	output, _ := s.Repository.GetUserData(ctx.Request().Context(), repository.UserInput{
		PhoneNumber: params.PhoneNumber,
	})

	if err := bcrypt.CompareHashAndPassword([]byte(output.Password), []byte(params.Password)); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid phone number or password")
	}

	// Create token
	claims := &Claims{
		PhoneNumber:    params.PhoneNumber,
		Name:           output.Name,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("sawit-pro-digital"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to generate token")
	}

	err = s.Repository.Logged(ctx.Request().Context(), params.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed login to account")
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

// PostSignup implements generated.ServerInterface.
func (s *Server) PostSignup(ctx echo.Context, params generated.PostSignupParams) error {
	isValid, errValMsg := validateInput(params.PhoneNumber, params.FullName, params.Password)
	if !isValid {
		return echo.NewHTTPError(http.StatusBadRequest, errValMsg)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password ")
	}

	user := repository.UserInput{
		PhoneNumber: params.PhoneNumber,
		Password:    hashedPassword,
		FullName:    params.FullName,
	}
	output, err := s.Repository.SignUp(ctx.Request().Context(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Account already exists")
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Successfully signed up",
		"ID":      strconv.Itoa(output.ID),
	})
}

func validateInput(phoneNumber, fullName, password string) (bool, string) {
	var errorMsgs []string

	if !validatePhoneNumber(phoneNumber) {
		errorMsgs = append(errorMsgs, "Phone number must be at least 10 characters, start with +62")
	}

	if !validateFullName(fullName) {
		errorMsgs = append(errorMsgs, "Full name must be between 3 and 60 characters")
	}

	if !validatePassword(password) {
		errorMsgs = append(errorMsgs, "Password must be between 6 and 64 characters, contain at least 1 uppercase letter, 1 number, and 1 special character")
	}

	if len(errorMsgs) > 0 {
		errorMsg := strings.Join(errorMsgs, ", ")
		return false, errorMsg
	}

	return true, ""
}

func validatePhoneNumber(phoneNumber string) bool {
	if !strings.HasPrefix(phoneNumber, "+62") {
		log.Println("++")
		return false
	}

	numberWithoutPrefix := strings.TrimPrefix(phoneNumber, "+62")

	for _, char := range numberWithoutPrefix {
		if char < '0' || char > '9' {
			log.Println("angg")
			return false
		}
	}

	if len(numberWithoutPrefix) < 9 || len(numberWithoutPrefix) > 12 {
		log.Println("preff")
		return false
	}

	return true
}

func validateFullName(fullName string) bool {
	if len(fullName) < 3 || len(fullName) > 60 {
		return false
	}

	return true
}

func validatePassword(password string) bool {
	if len(password) < 6 || len(password) > 64 {
		return false
	}

	hasUppercase := false
	hasNumber := false
	hasSpecial := false
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUppercase = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUppercase && hasNumber && hasSpecial
}
