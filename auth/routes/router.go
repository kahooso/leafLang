package routes

import (
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine) {

	r.SetFuncMap(template.FuncMap{
		"toLower": strings.ToLower,
	})

	RegisterPublicRoutes(r)
	RegisterAPIRoutes(r)
	RegisterProtectedRoutes(r)
}
