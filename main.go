package main

import (
	"com/go-jwt/config"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var cfg = config.GetConfig()

func CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()

	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized: " + err.Error()))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}

func GetJwt(w http.ResponseWriter, r *http.Request) {
	if r.Header["Access-Token"] != nil {
		if r.Header["Access-Token"][0] != cfg.JWTApiKey {
			return
		} else {
			token, err := CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprint(w, token)
		}
	}
}

func HelloHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	http.Handle("/hello", ValidateJWT(HelloHandle))
	http.HandleFunc("/jwt", GetJwt)

	http.ListenAndServe(":8080", nil)
}
