package sess

import (
	"fmt"

	guuid "github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Init Session init //session cookie 的 secrect key
func Init() echo.MiddlewareFunc {
	return session.Middleware(sessions.NewCookieStore([]byte("!@#$!%^$&^*DASFASGHJTRURWR")))
}

// Handler middleware adds Session
// cookie session to every request.
func Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 0,
			// MaxAge=0 means no 'Max-Age' attribute specified.
			// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
			// MaxAge>0 means Max-Age attribute present and given in seconds.
			HttpOnly: true,
		}
		//fmt.Println(c.Request().Method)
		//判斷有沒有 uuid
		if _, ok := sess.Values["SessionID"]; ok {
			//do something here
			//ssid = fmt.Sprintf("%s", ssid)
			//fmt.Printf("--has--ssid--%s---\r\n", ssid)
		} else {
			// /someuri 排除
			//if strings.Contains(c.Request().RequestURI, "/someuri") != true {
			uuid := guuid.New()
			//fmt.Printf("github.com/google/uuid:         %s\n\r\n", uuid.String())
			sess.Values["SessionID"] = fmt.Sprintf("%s", uuid.String())
			//fmt.Printf("RequestURI:         %s\n\r\n", c.Request().RequestURI)
			//}
		}
		sess.Save(c.Request(), c.Response())
		c.Response().Header().Set(echo.HeaderServer, "Echo/4.1")
		return next(c)
	}
}

// GetSession return sess
func GetSession(c echo.Context, key string) (interface{}, bool) {
	sess, _ := session.Get("session", c)
	var result interface{}
	if result, ok := sess.Values[key]; ok {
		return result, true
	}
	return result, false
}
