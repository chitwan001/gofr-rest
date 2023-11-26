package userRoutes

import (
	"github.com/golang-jwt/jwt"
	"gofr.dev/pkg/errors"
	"gofr.dev/pkg/gofr"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("my_secret_key")

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	HashPass string `json:"hash_pass"`
}

type UserRequestData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginBody login body
type LoginBody struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

// Claims : type for jwt body
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthMiddlewareBody struct {
	Username string `json:"userEmail"`
}

// HashPassword returns hashed password generated from the password passed as argument to it
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// RegisterUser creates a new user if user does not already exist in the database/**
func RegisterUser(ctx *gofr.Context) (interface{}, error) {
	var user UserRequestData
	err := ctx.Bind(&user)
	if err != nil {
		return nil, err
	}

	//check if any user already exists with the same email
	rows, queryErr := ctx.DB().QueryContext(ctx, "SELECT * FROM users WHERE email=?", user.Email)

	if queryErr != nil {
		return nil, queryErr
	}

	count := 0
	for rows.Next() {
		count++
		break
	}
	if count != 0 {
		return nil, &errors.Response{
			StatusCode: http.StatusForbidden,
			Reason:     "User already exists",
			ResourceID: user.Email,
		}
	} else {
		hashPassword, err := HashPassword(user.Pass)
		if err != nil {
			return nil, err
		}
		_, err = ctx.DB().ExecContext(ctx, "INSERT INTO users (name,email,hash_pass) VALUES (?,?,?)", user.Name, user.Email, hashPassword)

		//now login
		expirationTime := time.Now().Add(24 * 30 * 5 * time.Hour).Unix()
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: user.Email,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			return nil, err
		}
		success := make(map[string]string)
		success["description"] = "User successfully created and logged in!"
		success["token"] = tokenString
		success["statusCode"] = strconv.Itoa(http.StatusOK)
		return success, nil
	}
}

func LoginUser(ctx *gofr.Context) (interface{}, error) {
	var loginBody LoginBody
	err := ctx.Bind(&loginBody)
	if err != nil {
		return nil, err
	}

	//check if any user already exists with the same email
	rows, queryErr := ctx.DB().QueryContext(ctx, "SELECT * FROM users WHERE email=?", loginBody.Email)

	if queryErr != nil {
		return nil, queryErr
	}

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.HashPass); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, &errors.Response{
			StatusCode: http.StatusForbidden,
			Reason:     "User does not exist",
			ResourceID: loginBody.Email,
		}
	} else {
		//user exists check for password

		var user = users[0]

		compareErr := bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(loginBody.Pass))
		if compareErr != nil {
			//password does not match
			return nil, &errors.Response{
				StatusCode: http.StatusUnauthorized,
				Reason:     "Password does not matches",
				ResourceID: user.Email,
			}
		} else {
			//password matches then generate jwt and send
			//expirationTime is of 5 months
			expirationTime := time.Now().Add(24 * 30 * 5 * time.Hour).Unix()
			// Create the JWT claims, which includes the username and expiry time
			claims := &Claims{
				Username: loginBody.Email,
				StandardClaims: jwt.StandardClaims{
					// In JWT, the expiry time is expressed as unix milliseconds
					ExpiresAt: expirationTime,
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				// If there is an error in creating the JWT return an internal server error
				return nil, err
			}
			response := make(map[string]string)
			response["token"] = tokenString
			response["description"] = "User created and logged in successfully"
			return response, nil
		}
	}

}

func Me(ctx *gofr.Context) (interface{}, error) {
	email := ctx.Value("userEmail")

	rows, queryErr := ctx.DB().QueryContext(ctx, "SELECT name,email FROM users WHERE email=? LIMIT 1", email)

	if queryErr != nil {
		return nil, queryErr
	}

	var users []UserResponse
	for rows.Next() {
		var user UserResponse
		if err := rows.Scan(&user.Name, &user.Email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	response := map[string]interface{}{
		"me": users[0],
	}

	return response, nil
}

func DeleteUser(ctx *gofr.Context) (interface{}, error) {
	email := ctx.Value("userEmail")
	_, err := ctx.DB().ExecContext(ctx, "DELETE FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	//as delete request must not return anything but 204 No Content
	return nil, nil
}
