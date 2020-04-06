package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func SignUpResponse(c echo.Context, userId string) error {
	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}{
		Message: "Account successfully created",
		User: User{
			UserId:   userId,
			Nickname: userId,
		},
	})
}

func UserResponse(c echo.Context, user User) error {
	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}{
		Message: "User details by user_id",
		User: user,
	})
}

func UpdateResponse(c echo.Context, nickname, comment string) error {
	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		Recipe  User   `json:"user"`
	}{
		Message: "User successfully updated",
		Recipe: User{
			Nickname: nickname,
			Comment:  comment,
		},
	})
}

func DeleteResponse(c echo.Context) error {
	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Account and user successfully removed",
	})
}

func ErrorResponse(c echo.Context, code int, message, cause string) error {
	return c.JSON(code, struct {
		Message string `json:"message"`
		Cause   string `json:"cause,omitempty"`
	}{
		Message: message,
		Cause: cause,
	})
}

func AccountCreationError(c echo.Context, cause string) error {
	return ErrorResponse(c, http.StatusBadRequest, "Account creation failed", cause)
}

func AuthenticationError(c echo.Context) error {
	return ErrorResponse(c, http.StatusUnauthorized, "Authentication Faild", "")
}

func NoUserFoundError(c echo.Context) error {
	return ErrorResponse(c, http.StatusBadRequest, "No user found", "")
}

func NoPermissionError(c echo.Context) error {
	return ErrorResponse(c, http.StatusForbidden, "No Permission for Update", "")
}

func UserUpdationError(c echo.Context, cause string) error {
	return ErrorResponse(c, http.StatusBadRequest, "User updation failed", cause)
}

func DBConnectionError(c echo.Context, err string) error {
	return ErrorResponse(c, http.StatusInternalServerError, "Database connection failed", err)
}
