package routes

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func Set(r *gin.Engine) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	r.Static("/static", filepath.Join(cwd, "static"))
	r.LoadHTMLGlob(filepath.Join(cwd, "templates", "*"))

	r.SetFuncMap(template.FuncMap{
		"toLower": strings.ToLower,
	})

	RegisterPublicRoutes(r)
	RegisterAPIRoutes(r)
	RegisterProtectedRoutes(r)
}
