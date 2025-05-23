package controller

import (
	"net/http"
	"web/src/db"
	"web/src/lib"

	"github.com/gorilla/sessions"              // For sessions.Options
	"github.com/labstack/echo-contrib/session" // Ensure this is imported
	"github.com/labstack/echo/v4"
)

// ShowLoginForm renders the login page
func ShowLoginForm(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", lib.PyDict{
		"title": "Login",
	})
}

// Login handles the POST request from the login form
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	adminUser := db.AdminUser{}
	DB := db.Manager()

	// Find the user by username
	if err := DB.Where("username = ?", username).First(&adminUser).Error; err != nil {
		// User not found or other database error
		return c.Render(http.StatusOK, "login.html", lib.PyDict{
			"title": "Login",
			"error": "Invalid username or password 0",
		})
	}

	// Check the password
	if !adminUser.CheckPassword(password) {
		return c.Render(http.StatusOK, "login.html", lib.PyDict{
			"title": "Login",
			"error": "Invalid username or password",
		})
	}

	// Password is correct, create a session
	sess, _ := session.Get("session", c) // Use the echo-contrib session directly
	sess.Options = &sessions.Options{    // Set options before saving
		Path:     "/",
		MaxAge:   86400 * 7, // Example: 7 days
		HttpOnly: true,
		//Secure:   c.Request().IsTLS(), // Set Secure based on TLS
		// SameSite: http.SameSiteLaxMode, // Recommended for most cases
	}
	sess.Values["user_id"] = adminUser.ID
	sess.Values["username"] = adminUser.Username
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Error("Failed to save session:", err)
		return c.Render(http.StatusOK, "login.html", lib.PyDict{
			"title": "Login",
			"error": "Login failed. Please try again.",
		})
	}

	// Redirect to a protected area, e.g., home page or an admin dashboard
	return c.Redirect(http.StatusFound, "/")
}

// Logout clears the session and redirects to the login page
func Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	// Clear session values
	sess.Values["user_id"] = nil
	sess.Values["username"] = nil
	// Setting MaxAge to -1 deletes the cookie
	sess.Options = &sessions.Options{Path: "/", MaxAge: -1, HttpOnly: true}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		c.Logger().Error("Failed to save session on logout:", err)
		// Even if saving fails, attempt to redirect
	}
	return c.Redirect(http.StatusFound, "/login/")
}

// SetupAdminUser creates an initial admin user.
// IMPORTANT: This should be removed or protected after initial setup.
func SetupAdminUser(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "" || password == "" {
		return c.String(http.StatusBadRequest, "Please provide 'username' and 'password' query parameters.")
	}

	DB := db.Manager()

	var existingUser db.AdminUser
	if err := DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return c.String(http.StatusConflict, "Admin user with this username already exists.")
	}

	adminUser := db.AdminUser{
		Username:     username,
		PasswordHash: password, // The BeforeSave hook in AdminUser model will hash this
	}

	result := DB.Create(&adminUser)
	if result.Error != nil {
		c.Logger().Error("Failed to create admin user:", result.Error)
		return c.String(http.StatusInternalServerError, "Failed to create admin user: "+result.Error.Error())
	}

	return c.String(http.StatusOK, "Admin user '"+username+"' created successfully. Please remove or protect the /setup-admin route.")
}

// AuthMiddleware checks if a user is logged in
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		userID, ok := sess.Values["user_id"].(uint) // Ensure type assertion is correct for ID type
		if !ok || userID == 0 {
			// User not logged in, redirect to login page
			return c.Redirect(http.StatusFound, "/login/")
		}
		return next(c)
	}
}
