package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

/*
type Claims struct {
	username string
	jwt.RegisteredClaims
}
*/

func GenerateJWT(user string) (string, error) {
	/*
		expires_at := time.Now().Add(15 * time.Minute)

		claims := &Claims{
			username: user,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expires_at),
			},
		}
	*/
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token := jwt.New(jwt.SigningMethodHS256)
	singingKey := []byte(os.Getenv("SIGNINGKEY"))

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(15 * time.Minute)
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

		singingKey := []byte(os.Getenv("SIGNINGKEY"))

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
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized access: ", err)
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

	signingkey := []byte(os.Getenv("SIGNINGKEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", "There was an error in parsing token.")
		}
		return signingkey, nil
	})
	if err != nil {
		return "", err
	}
	// ... error handling

	// do something with decoded claims
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error decoding user info from token")
	}

	username := claim["username"].(string)

	return username, nil
}
