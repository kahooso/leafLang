package routes

import (
	"auth/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(r *gin.Engine) {
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.Use(middleware.ErrorMiddleware())
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/signin")
	})
	r.GET("/signin", func(c *gin.Context) {
		c.HTML(200, "signin.html", nil)
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})
	r.GET("/info", func(c *gin.Context) {
		c.HTML(200, "info.html", nil)
	})
	r.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "contact.html", nil)
	})
	r.GET("/faq", func(c *gin.Context) {
		c.HTML(200, "faq.html", nil)
	})
	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Redirect(http.StatusFound, "/signin")
	})
}
