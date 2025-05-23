package router

import (
	"fmt"
	"io"
	"os"

	//"fmt"

	"web/src/controller"
	"web/src/tpl" // template

	//"text/template" //text template
	"net/http"
	echopprof "web/src/pprof"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware" //v4 use v4 middleware
)

func GetIP(c echo.Context) string {

	// Open the log file in append mode
	logFile, err := os.OpenFile("log/ip.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err) // Handle error as appropriate for your application
	}
	defer logFile.Close()

	// Create a multiwriter to write to both the file and standard output
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Log all request headers to the file
	for name, headers := range c.Request().Header {
		for _, h := range headers {
			// Append header information to the log file
			_, err := fmt.Fprintf(multiWriter, "Header %v: %v\n", name, h)
			if err != nil {
				panic(err)
			}
		}
	}

	ip := ""

	// Prioritize headers in order: X-Forwarded-For, X-Real-IP, Sec-User-Ip
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "Sec-User-Ip"} {
		ip = c.Request().Header.Get(header)
		if ip != "" {
			c.Logger().Info("Using header: %s, IP: %s", header, ip)
			return ip
		}
	}

	if ip != "" {
		// If none of the preferred headers are found, fallback to c.RealIP()
		ip = c.RealIP()
	}

	c.Logger().Info("Using c.RealIP(), IP: %s", ip)
	return ip

}

// Init Router
func Init() *echo.Echo {

	e := echo.New()
	file, err := os.OpenFile("log/ip.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	e.Logger.SetOutput(file)
	e.Pre(middleware.AddTrailingSlash()) //request uri 最後面自動加 /

	t := tpl.Init()

	e.Renderer = t
	/*
		https://echo.labstack.com/guide/routing
		I would like to apply the same middlewares for all /admin/* URLs. Is it possible ?
		g := e.Group("/admin", <your-middleware>)
		g.GET("/secured", <your-handler>)
		Now route /admin/secured will execute your-middleware.
		https://github.com/labstack/echo/issues/613
	*/
	e.GET("/hello/:page/*", controller.Demo)

	e.GET("/hello/*", controller.Demo)

	// Group for authenticated routes
	adminGroup := e.Group("")                 // Apply to root or a specific path like /admin
	adminGroup.Use(controller.AuthMiddleware) // Use the AuthMiddleware

	adminGroup.GET("/", func(c echo.Context) error {
		// This route is now protected
		// You can get user from session if needed:
		// sess, _ := session.Get("session", c) // session.Get is from "github.com/labstack/echo-contrib/session"
		// For this to work, session middleware must be configured in main.go:
		// e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret-key"))))
		// We assume it's configured. Let's try to get username from session.
		// The AuthMiddleware should ensure session exists and has user_id.
		// We also stored "username" in the Login function.

		_sess, err_sess := session.Get("session", c)
		if err_sess != nil {
			// Handle error, maybe redirect to login or show error
			// This case should ideally be caught by AuthMiddleware redirecting if session is invalid
			return c.String(http.StatusInternalServerError, "Session error")
		}

		username, ok := _sess.Values["username"].(string)
		if !ok || username == "" {
			// This case should also ideally be caught by AuthMiddleware
			// If username is not in session, redirect or error
			return c.Redirect(http.StatusFound, "/login/")
		}

		userIP := GetIP(c)
		output := "Hello, " + username + "! Your IP is " + userIP
		return c.String(http.StatusOK, output)
	})

	e.Any("/user/:name", func(c echo.Context) error {
		name := c.Param("name")
		output := fmt.Sprintf("Hello, %s!", name)
		return c.String(http.StatusOK, output)
	})

	e.POST("/api/addTask/", controller.AddTask)

	//tasks
	e.GET("/api/tasks/:taskIDoruserID/*", controller.Tasks)
	e.GET("/api/tasks/*", controller.Tasks)

	//records
	e.GET("/api/records/:userID/*", controller.Records)
	e.GET("/api/records/*", controller.Records)

	//userPointBalance balance
	e.GET("/api/balance/:userID/*", controller.Users)
	e.GET("/api/balance/*", controller.Users)

	e.GET("/login/", controller.ShowLoginForm) // Trailing slash for consistency
	e.POST("/login/", controller.Login)
	///success/ return json
	e.GET("/success/", func(c echo.Context) error {
		output := fmt.Sprintf("Hello, Success!")
		return c.String(http.StatusOK, output)
	})
	//e.GET("/setup-admin/", controller.SetupAdminUser) // Add trailing slash
	e.GET("/logout/", controller.Logout) // Add trailing slash

	//所有 public folder 的文件都可以被 access ,例如 public/robots.txt -> http://yourhost/robots.txt
	e.Static("/", "public")

	//print payload
	e.Any("/api/print_all/", controller.PrintPayload)

	// EX:go tool pprof http://192.168.119.128/
	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	echopprof.Wrap(e)

	return e
}
