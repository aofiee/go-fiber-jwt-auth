package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/aofiee/diablos/diablosutils"
	"github.com/aofiee/diablos/types"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	utils "github.com/gofiber/fiber/v2/utils"
)

func TestSetup(t *testing.T) {
	t.Run("FAIL_INIT_CONFIG", func(t *testing.T) {
		var err error
		config, err = diablosutils.LoadConfig("../../")
		if err != nil {
			utils.AssertEqual(t, true, err)
		}
	})
	t.Run("SUCCESS_INIT_CONFIG", func(t *testing.T) {
		var err error
		config, err = diablosutils.LoadConfig("../")
		if err != nil {
			utils.AssertEqual(t, nil, err)
		}
	})
}

func Test_Auth(t *testing.T) {
	app := fiber.New()
	app.Post("/login", Auth)
	var params types.Login
	t.Run("FAIL_POST_JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
	})
	t.Run("FAIL_POST_USERNAME_PASSWORD", func(t *testing.T) {
		params.Username = "aofiee666"
		params.Password = "password"
		req := httptest.NewRequest("POST", "/login", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
	})
	params.Username = "aofiee"
	params.Password = "password"
	t.Run("FAIL_REDIS_CONNECTION", func(t *testing.T) {
		data, _ := json.Marshal(&params)
		payload := bytes.NewReader(data)
		req := httptest.NewRequest("POST", "/login", payload)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 500, resp.StatusCode, "Status code")
	})
}
func VerifyAccessToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(config.AccessKey), nil
}

/*
func Test_Profile(t *testing.T) {
	t.Run("CREATE_TOKEN", func(t *testing.T) {
		app := fiber.New()

		token, err := createToken("sawatdeekub")
		if err != nil {
			utils.AssertEqual(t, nil, err)
		}
		extoken, _ := jwt.Parse(token.Token.AccessToken, VerifyAccessToken)
		utils.AssertEqual(t, true, extoken.Valid)

		claims, ok := extoken.Claims.(jwt.MapClaims)
		utils.AssertEqual(t, true, ok)

		uid, ok := claims["sub"].(string)
		utils.AssertEqual(t, "sawatdeekub", uid)
		utils.AssertEqual(t, true, ok)

		// var ctx fiber.Ctx
		// ctx.Locals("user", extoken)

		bearer := "Bearer " + token.Token.AccessToken
		app.Get("/profile", func(c *fiber.Ctx) error {
			return Profile(c)
		})

		req := httptest.NewRequest("GET", "/profile", nil)
		req.Header.Set("Authorization", bearer)
		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
	})
}
*/
func Test_deleteAuthFromRedis(t *testing.T) {
	var ctx *fiber.Ctx
	d, err := deleteAuthFromRedis(ctx, "hello")
	utils.AssertEqual(t, "", d, "d")
	utils.AssertEqual(t, "dial tcp: lookup redis: no such host", err.Error(), "err")
}

func Test_AuthError(t *testing.T) {
	a, err := FetchAuth("xxxx")
	utils.AssertEqual(t, "dial tcp: lookup redis: no such host", err.Error(), "err")
	utils.AssertEqual(t, "", a, "a")
}
