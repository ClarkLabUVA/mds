package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"

	"strings"
)

var JWTSECRET = "orstestsecret"

type ORSClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

// Parses Token from Header and adds a User Object to the passed context
// Simply Passes if
func JWTAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var u User

		// get the JWT from the header
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// parse the jwt
		token, err := jwt.ParseWithClaims(tokenString, &ORSClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSECRET), nil
		})

		if err != nil {
			h.ServeHTTP(w, r)
		}

		if claims, ok := token.Claims.(*ORSClaims); ok && token.Valid {
			u.ID = claims.Subject
			u.Type = "Person"
			u.Name = claims.Name
			u.Email = claims.Email

			// set context variable for request
			ctx := context.WithValue(r.Context(), "User", u)
			// pass to next handler
			h.ServeHTTP(w, r.WithContext(ctx))
		} else {
			h.ServeHTTP(w, r)
		}

	})
}

func JWTAuthRequired(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
