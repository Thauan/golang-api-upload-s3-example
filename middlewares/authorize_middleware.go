package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Thauan/golang-api-upload-s3-example/handlers"
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("Invalid")
				}

				return []byte(handlers.GetEnvWithKey("SECRET_KEY")), nil
			})

			if token == nil {
				data, _ := json.Marshal("Invalid Token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(data)
				return
			}

			if err != nil {
				data, _ := json.Marshal(err.Error())

				w.WriteHeader(http.StatusBadGateway)
				w.Write(data)
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {

			data, _ := json.Marshal("No Authorization Token provided")

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(data)

		}
	})
}
