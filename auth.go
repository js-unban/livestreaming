package main

import (
	"example/config"
	"example/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LivestreamingCustomClaims struct {
	UserId uint `json:"userId"`
	jwt.RegisteredClaims
}

const AUTHORIZATION_HEADER = "Authorization"
const COOKIE_SESSION_PARAM = "sess"

type LoginRequest struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func getUserPasswordFromForm(c *gin.Context) (LoginRequest, error) {
	var loginRequest LoginRequest
	if err := c.ShouldBind(&loginRequest); err != nil {
		return loginRequest, err
	} else {
		return loginRequest, nil
	}

}

func generateTokenAndRespond(c *gin.Context, userId uint, successRedirect string, errorTemplate string) error {
	token, err := GenerateUserTokenFromId(uint32(userId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, errorTemplate, gin.H{"message": "user not found"})
		return err
	}
	c.SetCookie(COOKIE_SESSION_PARAM, token, 3600, "/", "localhost", true, true)
	c.Redirect(http.StatusTemporaryRedirect, successRedirect)
	return nil
}

func SetUserFromAuthHMAC(c *gin.Context) {
	var claims LivestreamingCustomClaims
	authorization := strings.TrimSpace(strings.Replace(c.Request.Header.Get(AUTHORIZATION_HEADER), "Bearer", "", 1))
	token, err := jwt.ParseWithClaims(authorization, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method from received token")
		}
		return []byte(models.SIGNING_KEY), nil
	})

	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("bad token %s %s", authorization, err.Error()))
		c.Abort()
		return
	}

	if !token.Valid {
		sendErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	c.Set(GIN_USER_ID_KEY, claims.UserId)
	c.Next()

}

func GenerateUserTokenFromId(id uint32) (string, error) {
	cfg := config.GetConfig()
	signingKey := []byte(models.SIGNING_KEY)
	futureTime := time.Now().Add(time.Hour * time.Duration(cfg.TOKEN_DURATION_HOURS))

	if user, ok := models.Users.FindUserById(uint(id)); !ok {
		msg := "username not found"
		log.Debug(msg)
		return "", errors.New(msg)
	} else {

		claims := LivestreamingCustomClaims{
			uint(user.ID),
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(futureTime),
				Issuer:    cfg.ISSUER,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if ss, err := token.SignedString(signingKey); err != nil {
			wErr := errors.Wrap(err, "error signingKey")
			log.Error(wErr.Error())
			return "", wErr
		} else {
			return ss, nil
		}
	}
}

func redirectToLogin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/login")
	c.Abort()
}
