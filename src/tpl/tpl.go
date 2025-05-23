package tpl

import (

	//"text/template" //text template
	"html/template" // html template
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

//Template function
type Template struct {
	templates *template.Template
}

//Render Template function
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//Init function
func Init() *Template {
	//template custom function
	funcMap := template.FuncMap{
		"strupper": func(str string) string { return strings.ToUpper(str) },
		"hasPermission": func(feature int) bool {
			if feature == 1 {
				return true
			}
			return false
		},
		"htmlSafe": func(html string) template.HTML {
			return template.HTML(html)
		},
	}

	t1 := template.New("All")
	t1 = t1.Funcs(funcMap)
	//resource/views/*.html is the template path
	t1, _ = t1.ParseGlob("resource/views/*.html")
	t := &Template{templates: t1}

	return t
}
