package main

import (
	"context"
	"github.com/golang-jwt/jwt"
	"gofr.dev/pkg/gofr"
	"golang-form-services/userRoutes"
	"net/http"
	"strings"
)

func baseHandler(ctx *gofr.Context) (interface{}, error) {
	resp := make(map[string]string)
	resp["message"] = "Status Created"
	return resp, nil
}

var jwtKey = []byte("my_secret_key")

// Claims : type for jwt body
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func authMiddleware() func(handler http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			// if request URI is /users/create or /users/login or /(base url) then no need for authentication check
			requestURI := r.RequestURI
			if requestURI == "/user/create" || requestURI == "/user/login" || requestURI == "/" {
				//no need for authentication
				// sends the request to the next middleware/request-handler
				inner.ServeHTTP(w, r)
			} else {
				//user has to be authenticated here
				authHeader := r.Header["Authorization"]
				if authHeader == nil {
					http.Error(w, "Not authorized. Please login!", http.StatusUnauthorized)
					return
				}
				if authHeader[0] == "" {
					http.Error(w, "Not authorized. Please login!", http.StatusUnauthorized)
					return
				}
				tokenStr := strings.Split(authHeader[0], " ")
				if len(tokenStr) == 1 {
					http.Error(w, "Not authorized. Please login!", http.StatusUnauthorized)
					return
				}

				claims := &Claims{}

				// Parse the JWT string and store the result in `claims`.
				// Note that we are passing the key in this method as well. This method will return an error
				// if the token is invalid (if it has expired according to the expiry time we set on sign in),
				// or if the signature does not match
				tkn, err := jwt.ParseWithClaims(tokenStr[1], claims, func(token *jwt.Token) (any, error) {
					return jwtKey, nil
				})
				if err != nil || !tkn.Valid {
					http.Error(w, "You are unauthorized. Please login again!", http.StatusUnauthorized)
					return
				}

				//token is valid and not expired hence a valid and authorized user
				newContext := context.WithValue(r.Context(), "userEmail", claims.Username)
				newR := r.Clone(newContext)

				// sends the request to the next middleware/request-handler
				inner.ServeHTTP(w, newR)
			}
		})
	}
}

func main() {
	// initialise gofr object
	app := gofr.New()
	app.Server.UseMiddleware(authMiddleware())

	// register route greet
	app.GET("/", baseHandler)

	//user routes
	app.POST("/user/create", userRoutes.RegisterUser)
	app.POST("/user/login", userRoutes.LoginUser)
	app.GET("/user/me", userRoutes.Me)

	// Starts the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Start()
}
