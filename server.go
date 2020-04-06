package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {
	e := echo.New()
	u := UserModel{}
	h := Handler{&u}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/signup", h.SignUp)
	e.GET("users/:user_id", h.GetUser)
	e.PATCH("users/:user_id", h.PatchUser)
	e.POST("/close", h.CloseUser)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
