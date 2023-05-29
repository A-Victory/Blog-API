package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

func GenerateJWT(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	singingKey := []byte("secretkey")
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user
	claims["authorized"] = true
	claims["Exp"] = time.Now().Add(15 * time.Minute)
	tokenString, err := token.SignedString(singingKey)
	if err != nil {
		fmt.Println("Error signing token: ", err)
	}

	return tokenString, nil

}

func Verify(endpoint func(w http.ResponseWriter, r *http.Request, _ httprouter.Params)) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Header.Get("Token") == "" {
			_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
			return
		}

		singingKey := []byte("secretkey")

		tokenstring := r.Header["Token"][0]
		if tokenstring == "" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "User not authorized please login!")
			return
		}

		token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%v", "There was an error in parsing token.")
			}
			return singingKey, nil
		})

		if err != nil {
			fmt.Println("Error while parsing token: ", err)
			return
		}

		if token.Valid {
			endpoint(w, r, httprouter.Params{})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}
	})

}

func GetUser(r *http.Request) (string, error) {
	tokenString := r.Header["Token"][0]
	if tokenString == "" {
		return "", errors.New("token string is empty")
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("<YOUR VERIFICATION KEY>"), nil
	})
	if err != nil {
		return "Error while parsing token: ", err
	}
	// ... error handling

	// do something with decoded claims
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)

	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		username := claims["username"].(string)
		return username, nil
	}

	cusErr := fmt.Errorf("Could not extract user from token")

	return "", cusErr
}
