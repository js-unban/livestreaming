package main

import (
	"example/models"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func ValidateSessionFromCookie(c *gin.Context) {
	var claims LivestreamingCustomClaims
	cookie, err := c.Cookie(COOKIE_SESSION_PARAM)
	if err != nil {
		redirectToLogin(c)
	}

	authorization := strings.TrimSpace(cookie)
	token, err := jwt.ParseWithClaims(authorization, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method from received token")
		}
		return []byte(models.SIGNING_KEY), nil
	})

	if err != nil {
		// bad input for token
		redirectToLogin(c)
		return
	}

	if !token.Valid {
		//not valid token
		redirectToLogin(c)
		return
	}

	c.Set(GIN_USER_ID_KEY, claims.UserId)
	c.Next()

}

func streamVOD(c *gin.Context) {

	c.Stream(func(w io.Writer) bool {

		w.Write([]byte(STATIC_HEADER))
		w.Write([]byte(strconv.Itoa(counter)))
		w.Write([]byte("\n"))
		for i := counter; i < counter+DEFAULT_SEGMENT_COUNT; i++ {
			w.Write([]byte(fmt.Sprintf(SEGMENT_TEMPLATE, i)))

		}

		if counter >= maxCounter {
			//log.Println("send ending at", counter)
			w.Write([]byte(STATIC_ENDING))
		}

		counter += 1
		counter = counter % (maxCounter + 1)

		return false
	})

}

func register(c *gin.Context) {
	template := "register.tmpl"
	registrationRequest, err := getRegistrationParams(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, template, gin.H{"message": "bad params for registration, should receive name, username and password"})
		return
	}
	user, err := models.Users.CreateUser(&models.User{Name: registrationRequest.Name, Email: registrationRequest.Email, Password: registrationRequest.Password})
	if err != nil {
		var msg string
		if !errors.Is(err, models.UserAlreadyExists) {
			msg = "internal error creating user"
			log.Error(errors.Wrap(err, msg))
		} else {
			msg = err.Error()
		}
		c.HTML(http.StatusInternalServerError, template, gin.H{"message": msg})
		return
	}
	generateTokenAndRespond(c, user.ID, "/load_player", template)

}

func login(c *gin.Context) {
	template := "login.tmpl"
	loginRequest, err := getLoginParams(c)
	if err != nil {
		c.HTML(http.StatusBadRequest, template, gin.H{"message": "bad params for login, should receive username and password"})
		return
	}

	if ok := models.Users.ValidateUserPassword(loginRequest.Email, loginRequest.Password); !ok {
		c.HTML(http.StatusBadRequest, template, gin.H{"message": "invalid email/password"})
		return
	}

	queryRes, ok := models.Users.FindUserByEmail(loginRequest.Email)
	if !ok {
		// should never happen, previous method validates this as well but still handle
		c.HTML(http.StatusBadRequest, template, gin.H{"message": "user not found"})
		return
	}

	generateTokenAndRespond(c, queryRes.ID, "/load_player", template)

}

func logout(c *gin.Context) {
	c.SetCookie(COOKIE_SESSION_PARAM, "", 0, "/", "localhost", true, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
