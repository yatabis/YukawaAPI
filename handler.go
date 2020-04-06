package main

import (
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type Handler struct {
	user Model
}

func (handler *Handler) BasicAuth(c echo.Context) (userId string, ok bool) {
	userId, password, ok := c.Request().BasicAuth()
	pw := handler.user.FetchPassword(userId)
	if !ok || pw != password {
		return "", false
	}
	return
}

func (handler *Handler) SignUp(c echo.Context) error {
	user := new(struct{
		UserId   string `json:"user_id"  validate:"required,min=6,max=20,alphanum"`
		Password string `json:"password" validate:"required,min=8,max=20,printascii,excludes= "`
	})

	if err := c.Bind(user); err != nil {
		return AccountCreationError(c, err.Error())
	}

	field, tag := Validate(user)
	if tag == "required" {
		return AccountCreationError(c, "required user_id and password")
	}
	if field == "UserId" {
		if tag == "alphanum" {
			return AccountCreationError(c, "the pattern of user_id must be alphanumeric")
		} else {
			return AccountCreationError(c, "the length of user_id must be at least 6 and no more 20")
		}
	}
	if field == "Password" {
		if tag == "printascii" {
			return AccountCreationError(c, "the pattern of password must be ASCII characters")
		} else {
			return AccountCreationError(c, "the length of password must be at least 8 and no more 20")
		}
	}

	if err := handler.user.OpenDB(); err != nil {
		return DBConnectionError(c, err.Error())
	}
	defer handler.user.CloseDB()

	userId, _, _ := handler.user.FetchDetail(user.UserId)
	if userId != "" {
		return AccountCreationError(c, "already same user_id is used")
	}

	handler.user.New(user.UserId, user.Password)

	return SignUpResponse(c, user.UserId)
}

func (handler *Handler) GetUser(c echo.Context) error {
	if err := handler.user.OpenDB(); err != nil {
		return DBConnectionError(c, err.Error())
	}
	defer handler.user.CloseDB()

	if _, ok := handler.BasicAuth(c); !ok {
		return AuthenticationError(c)
	}

	userId, nickname, comment := handler.user.FetchDetail(c.Param("user_id"))
	if userId == "" {
		return NoUserFoundError(c)
	}

	if nickname == "" {
		nickname = userId
	}

	user := User{
		UserId:   userId,
		Nickname: nickname,
		Comment:  comment,
	}

	return UserResponse(c, user)
}

func (handler *Handler) PatchUser(c echo.Context) error {
	if err := handler.user.OpenDB(); err != nil {
		return DBConnectionError(c, err.Error())
	}
	defer handler.user.CloseDB()

	userId, ok := handler.BasicAuth(c)
	if !ok {
		return AuthenticationError(c)
	}

	if userId != c.Param("user_id") {
		return NoPermissionError(c)
	}

	_, nickname, comment := handler.user.FetchDetail(userId)

	user := new(struct{
		UserId   string `json:"user_id" validate:"isdefault"`
		Password string `json:"password" validate:"isdefault"`
		Nickname interface{} `json:"nickname" validate:"excludes_control,max=30"`
		Comment  interface{} `json:"comment" validate:"excludes_control,max=100"`
	})
	if err := c.Bind(&user); err != nil {
		return UserUpdationError(c, err.Error())
	}

	field, tag := Validate(user)
	if tag == "isdefault" {
		return UserUpdationError(c, "not updatable user_id and password")
	}

	if user.Nickname == nil && user.Comment == nil {
		return UserUpdationError(c, "required nickname or comment")
	}
	if user.Nickname == nil {
		user.Nickname = nickname
	}
	if user.Comment == nil {
		user.Comment = comment
	}
	field, tag = Validate(user)

	if field == "Nickname" {
		if tag == "max" {
			return UserUpdationError(c, "the length of nickname must be 30 or less")
		} else {
			return UserUpdationError(c, "nickname must not contain any control code")
		}
	}
	if field == "Comment" {
		if tag == "max" {
			return UserUpdationError(c, "the length of comment must be 100 or less")
		} else {
			return UserUpdationError(c, "comment must not contain any control code")
		}
	}

	handler.user.Update(userId, user.Nickname.(string), user.Comment.(string))
	if user.Nickname == "" {
		user.Nickname = userId
	}
	return UpdateResponse(c, user.Nickname.(string), user.Comment.(string))
}

func (handler *Handler) CloseUser(c echo.Context) error {
	if err := handler.user.OpenDB(); err != nil {
		return DBConnectionError(c, err.Error())
	}
	defer handler.user.CloseDB()

	userId, ok := handler.BasicAuth(c)
	if !ok {
		return AuthenticationError(c)
	}

	handler.user.Delete(userId)

	return DeleteResponse(c)
}
