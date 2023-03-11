package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const GIN_USER_ID_KEY = "userId"

func main() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	initMaxCounter()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:8080"}
	//r.Use(cors.New(config))
	r.GET("/", ValidateSessionFromCookie, func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/player")
	})
	addStreamControllers(r)
	addAuthControllers(r)

	r.Run()
}

func addStreamControllers(r *gin.Engine) {
	streamGroup := r.Group("/stream")
	{
		streamGroup.Use(ValidateSessionFromCookie)
		streamGroup.GET("/segment.m3u8", ValidateSessionFromCookie, streamVOD)
		streamGroup.Static("/segments", PLAYLIST_FOLDER+"segments")
	}
	r.GET("/player", ValidateSessionFromCookie, func(c *gin.Context) {
		c.HTML(http.StatusOK, "player.tmpl", gin.H{})
	})

	r.Static("/static", "./static")
	r.Any("/load_player", ValidateSessionFromCookie, func(c *gin.Context) {
		c.HTML(http.StatusOK, "load_player.tmpl", gin.H{})
	})
}

func addAuthControllers(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{})
	})
	r.POST("/login", login)

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{})
	})
	r.POST("/register", register)
	r.GET("/logout", logout)
}
