package main

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type userModelStub struct {
}

func (u *userModelStub) OpenDB() error { return nil }
func (u *userModelStub) CloseDB() {}
func (u *userModelStub) New(_, _ string) {}
func (u *userModelStub) FetchDetail(userId string) (string, string, string) {
	if userId == "TaroYamada" {
		return "TaroYamada", "たろー", "僕は元気です"
	}
	return "", "", ""
}
func (u *userModelStub) FetchPassword(_ string) string { return "test@password" }
func (u *userModelStub) Update(_, _, _ string) {}
func (u *userModelStub) Delete(_ string) {}

func TestSignup(t *testing.T) {
	e := echo.New()
	u := userModelStub{}
	h := Handler{&u}

	e.POST("/signup", h.SignUp)
	e.GET("users/:user_id", h.GetUser)
	e.PATCH("users/:user_id", h.PatchUser)
	e.POST("/close", h.CloseUser)

	tests := []struct {
		name string
		body string
		code int
		resp interface{}
	}{
		{
			"OK case",
			`{"user_id": "testUser", "password": "testPassWord"}`,
			http.StatusOK,
			`{"message":"Account successfully created","user":{"user_id":"testUser","nickname":"testUser"}}` + "\n",
		},
		{
			"no user_id",
			`{"password": "testPass"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"required user_id and password"}` + "\n",
		},
		{
			"short user_id",
			`{"user_id": "test", "password": "4wordUserId"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"the length of user_id must be at least 6 and no more 20"}` + "\n",
		},
		{
			"long user_id",
			`{"user_id": "testUserIdWith26Characters", "password": "26wordUserId"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"the length of user_id must be at least 6 and no more 20"}` + "\n",
		},
		{
			"user_id with space",
			`{"user_id": "test user", "password": "spaceuser"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"the pattern of user_id must be alphanumeric"}` + "\n",
		},
		{
			"no password",
			`{"user_id": "testId"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"required user_id and password"}` + "\n",
		},
		{
			"short password",
			`{"user_id": "pass6chrs", "password": "testpw"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"the length of password must be at least 8 and no more 20"}` + "\n",
		},
		{
			"long password",
			`{"user_id": "pass26chrs", "password": "abcdefghijklmnopqrstuvwxyz"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"the length of password must be at least 8 and no more 20"}` + "\n",
		},
		{
			"password with control code",
			`{"user_id": "control", "password": "pass\nword"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"the pattern of password must be ASCII characters"}` + "\n",
		},
		{
			"existing user_id",
			`{"user_id": "TaroYamada", "password": "password"}`,
			http.StatusBadRequest,
			`{"message":"Account creation failed","cause":"already same user_id is used"}` + "\n",
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/signup", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			if rec.Code != tt.code {
				t.Errorf("want = %d, got = %d", tt.code, rec.Code)
			}
			if got := rec.Body.String(); got != tt.resp.(string) {
				t.Errorf("want = %s, got = %s", tt.resp, got)
			}
		})
	}
}
