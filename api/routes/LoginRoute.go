package routes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aofiee/diablos/diablosutils"
	"github.com/aofiee/diablos/types"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	jwtware "github.com/gofiber/jwt/v2"
)

type (
	MsgLogin        types.Login
	MsgToken        types.Token
	MsgTokenDetail  types.TokenDetail
	MsgRefreshToken types.RefrestToken
	MsgJWTContext   types.JWTContext
)

var (
	rdConn *redis.Client
	//rbqDNS    = "amqp://" + os.Getenv("RB_USER") + ":" + os.Getenv("RB_PASSWORD") + "@" + os.Getenv("RB_HOST") + ":" + os.Getenv("RB_PORT") + "/"
	rbmqATExp = "900000"
	rbmqRTExp = "604800000"
	config    diablosutils.Config
)

func init() {
	var err error
	config, err = diablosutils.LoadConfig("../")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	rdAddr := config.RdHost + ":" + config.RdPort
	rdConn = redis.NewClient(&redis.Options{
		Addr:     rdAddr,
		Password: config.RdPassword,
		// DB:       0,
	})
}
func failOnError(c *fiber.Ctx, err error, msg string, status int) error {
	if err != nil {
		c.Status(status).JSON(fiber.Map{
			"msg":   msg,
			"error": err.Error(),
		})
	}
	return nil
}

func Auth(c *fiber.Ctx) error {
	var l MsgLogin
	uid := "144479bd-fcdc-4c9f-b116-f2a08807a4c3" //utils.UUID()
	err := c.BodyParser(&l)
	if err != nil {
		return failOnError(c, err, "cannot parse json", fiber.StatusBadRequest)
	}
	if l.Username != "aofiee" || l.Password != "password" {
		return failOnError(c, err, "Bad Credentials", fiber.StatusUnauthorized)
	}
	t, err := createToken(uid)
	if err != nil {
		return failOnError(c, err, "StatusForbidden", fiber.StatusForbidden)
	}
	err = storeJWTAuthToRedis(c, uid, t)
	if err != nil {
		return failOnError(c, err, "StatusInternalServerError", fiber.StatusInternalServerError)
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  t.Token.AccessToken,
		"refresh_token": t.Token.RefreshToken,
	})
	return nil
}

func generateTokenBy(uid string, rdKey string, ctx interface{}, signed string, expire int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	location, _ := time.LoadLocation("Asia/Bangkok")
	claims["iss"] = config.AppName
	claims["sub"] = uid
	claims["exp"] = expire
	claims["iat"] = time.Now().In(location).Unix()
	claims["context"] = ctx
	if ctx != nil {
		claims["access_uuid"] = rdKey
	} else {
		claims["refresh_uuid"] = rdKey
	}
	t, err := token.SignedString([]byte(signed))
	if err != nil {
		return "", err
	}
	return t, nil
}

func createToken(uid string) (*MsgTokenDetail, error) {
	tokenDetail := &MsgTokenDetail{}
	location, _ := time.LoadLocation("Asia/Bangkok")
	at := time.Now().In(location).Add(time.Minute * 15).Unix()
	rt := time.Now().In(location).Add(time.Hour * 24 * 7).Unix()
	tokenDetail.AccessTokenExp = at
	tokenDetail.RefreshTokenExp = rt
	tokenDetail.AccessUUid = utils.UUIDv4()
	tokenDetail.RefreshUUid = utils.UUIDv4()
	var err error
	/////mock data/////
	ctAccess := MsgJWTContext{}
	ctAccess.User = "aofiee666@gmail.com"
	ctAccess.DisplayName = "Khomkrid L."
	roles := []string{
		"admin",
		"report",
	}
	ctAccess.Roles = roles
	tokenDetail.Context = ctAccess
	/////mock data/////
	tokenDetail.Token.AccessToken, err = generateTokenBy(uid, tokenDetail.AccessUUid, ctAccess, config.AccessKey, tokenDetail.AccessTokenExp)
	if err != nil {
		return tokenDetail, err
	}
	tokenDetail.Token.RefreshToken, err = generateTokenBy(uid, tokenDetail.RefreshUUid, nil, config.RefreshKey, tokenDetail.RefreshTokenExp)
	if err != nil {
		return tokenDetail, err
	}
	return tokenDetail, nil
}

func deleteAuthFromRedis(c *fiber.Ctx, uid string) (string, error) {
	var ctx = context.Background()
	deleted, err := rdConn.GetDel(ctx, uid).Result()
	if err != nil {
		return "", err
	}
	return deleted, nil
}

func storeJWTAuthToRedis(c *fiber.Ctx, uid string, t *MsgTokenDetail) error {
	var err error
	var ctx = context.Background()
	location, _ := time.LoadLocation("Asia/Bangkok")
	atExp := time.Unix(t.AccessTokenExp, 0).In(location)
	rtExp := time.Unix(t.RefreshTokenExp, 0).In(location)
	now := time.Now().In(location)
	err = rdConn.Set(ctx, t.AccessUUid, uid, atExp.Sub(now)).Err()
	if err != nil {
		return err
	}
	err = rdConn.Set(ctx, t.RefreshUUid, uid, rtExp.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func AuthorizationRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: AuthSuccess,
		ErrorHandler:   AuthError,
		SigningKey:     []byte(config.AccessKey),
		SigningMethod:  "HS256",
		TokenLookup:    "header:Authorization",
		AuthScheme:     "Bearer",
	})
}

func AuthError(c *fiber.Ctx, e error) error {
	return failOnError(c, e, "Unauthorized", fiber.StatusUnauthorized)
}

func AuthSuccess(c *fiber.Ctx) error {
	c.Next()
	return nil
}

func FetchAuth(accessUUid string) (string, error) {
	var ctx = context.Background()
	uid, err := rdConn.Get(ctx, accessUUid).Result()
	if err != nil {
		return "", err
	}
	return uid, nil
}

func Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accessUUid := claims["access_uuid"].(string)
	context := claims["context"]
	uid, err := FetchAuth(accessUUid)
	if err != nil {
		return failOnError(c, err, "StatusUnauthorized", fiber.StatusUnauthorized)
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"sub":     uid,
		"context": context,
	})
	return nil
}

func Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accessUUid := claims["access_uuid"].(string)
	var err error
	_, err = deleteAuthFromRedis(c, accessUUid)
	if err != nil {
		return failOnError(c, err, "StatusUnauthorized", fiber.StatusUnauthorized)
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":   "Successfully logged out. Message : ",
		"error": nil,
	})
	return nil
}

func VerifyToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(config.RefreshKey), nil
}

func RefreshToken(c *fiber.Ctx) error {
	var rt MsgRefreshToken
	var err error
	err = c.BodyParser(&rt)
	if err != nil {
		return failOnError(c, err, "cannot recieve parameter", fiber.StatusBadRequest)
	}
	token, err := jwt.Parse(rt.Token, VerifyToken)
	if err != nil {
		return failOnError(c, err, "token signing error", fiber.StatusUnauthorized)
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return failOnError(c, err, "Refresh token expired", fiber.StatusUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, ok := claims["sub"].(string)
		if !ok {
			return failOnError(c, err, "token signing error", fiber.StatusNotFound)
		}
		refreshAccesss, ok := claims["refresh_uuid"].(string)
		if !ok {
			return failOnError(c, err, "token signing error", fiber.StatusNotFound)
		}
		_, err = deleteAuthFromRedis(c, refreshAccesss)
		if err != nil {
			return failOnError(c, err, "StatusForbidden", fiber.StatusForbidden)
		}
		t, err := createToken(uid)
		if err != nil {
			return failOnError(c, err, "StatusForbidden", fiber.StatusForbidden)
		}
		err = storeJWTAuthToRedis(c, uid, t)
		if err != nil {
			return failOnError(c, err, "StatusForbidden", fiber.StatusForbidden)
		}
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"access_token":  t.Token.AccessToken,
			"refresh_token": t.Token.RefreshToken,
		})
		return nil
	} else {
		return failOnError(c, err, "refresh expired", fiber.StatusUnauthorized)
	}
}
