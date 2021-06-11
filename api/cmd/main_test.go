package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/aofiee/diablos/types"

	// fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2"
	utils "github.com/gofiber/fiber/v2/utils"
)

type (
	TLogin types.Login
)

var (
	app *fiber.App
)

func TestMain(t *testing.T) {
	app = Setup()
	utils.AssertEqual(t, "*fiber.App", reflect.TypeOf(app).String(), "Setup()")
	t.Run("SUCCESS_ROOT", func(t *testing.T) {
		log.Println("app", reflect.TypeOf(app))
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
	})
	t.Run("FAIL_LOGIN", func(t *testing.T) {
		var params TLogin
		params.Username = "aofiee"
		params.Password = "password"
		data, _ := json.Marshal(&params)
		payload := bytes.NewReader(data)
		req := httptest.NewRequest("POST", "/login", payload)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 500, resp.StatusCode, "Status code")
	})
	t.Run("FAIL_PROFILE", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/profile", nil)
		req.Header.Set("Authorization", "Bearer xxx")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
	})
	t.Run("FAIL_LOGOUT", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/logout", nil)
		req.Header.Set("Authorization", "Bearer xxx")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
	})
	t.Run("FAIL_REFRESH_FORM", func(t *testing.T) {
		data := url.Values{}
		data.Set("refresh_token", "xxxx")
		payload := bytes.NewBufferString(data.Encode())
		req := httptest.NewRequest("POST", "/refresh", payload)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
	})
	t.Run("FAIL_REFRESH_JSON", func(t *testing.T) {
		params := map[string]interface{}{
			"refresh_token": "xxxx",
		}
		data, _ := json.Marshal(&params)
		payload := bytes.NewReader(data)
		req := httptest.NewRequest("POST", "/refresh", payload)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
	})
}
