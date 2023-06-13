package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/A-Victory/blog-API/auth"
	"github.com/A-Victory/blog-API/db"
	"github.com/A-Victory/blog-API/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UserController struct {
	Db *db.DbConn
	Va *auth.Validation
}

func NewUserController(d *db.DbConn, v *auth.Validation) *UserController {
	return &UserController{
		Db: d,
		Va: v,
	}
}

// Signup to create a new user account
func (uc UserController) Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	user.Id = primitive.NewObjectID()
	hashpswd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Fatal("Error generating password: ", err)
	}
	user.Password = string(hashpswd)

	err = uc.Va.ValidateUserInfo(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, _, err := uc.Db.GetUser(user)
	if username == user.Username {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Username already exists! Try again...")
		return
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			insert, err := uc.Db.CreateUser(user)
			if err != nil {
				http.Error(w, "Error creating accounting, please try again...", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(insert)

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusAlreadyReported)
		fmt.Fprintln(w, "Email already registered to an account! Please try again...")
		return
	}

	report := "Account successfully created!"
	json.NewEncoder(w).Encode(report)

}

// Login allows user to login in to existing account, creating a JWT token in the process.
func (uc UserController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	user := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	// Check the validity of information provided by the user in database
	if user.Email == "" {
		fmt.Fprintln(w, "Please input an email address!")
	}

	username, passwrd, err := uc.Db.GetUser(user)
	if err != nil {
		log.Println(err.Error())
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Email not registered, signup...", http.StatusInternalServerError)
			return
		}
		http.Error(w, "An error occured, try again...", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwrd), []byte(user.Password)); err != nil {
		fmt.Fprintln(w, "Incorrect password!!")
		return
	}

	token, err := auth.GenerateJWT(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Token", token)
	c := &http.Cookie{
		Name:     "user",
		Value:    user.Username,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	}
	http.SetCookie(w, c)

	report := "Login successful!"
	json.NewEncoder(w).Encode(report)
}

// UpdateInfo updates the user resources in the database.
func (uc UserController) UpdateInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	w.Header().Set("Content-Type", "application/json")
	username, err := auth.GetUser(r)
	json.NewEncoder(w).Encode(username)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	user := models.User{}
	//ctx := context.Background()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	err = uc.Db.UpdateUser(w, username, &user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Couldn't update user info!", http.StatusInternalServerError)
		return
	}
	// Return success if no error is encountered.]

	fmt.Fprintln(w, "user details successfully updated. ")

}

// Search allows a user search for other users. Returning information including the user's posts.
func (uc UserController) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	user := ps.ByName("name")
	username := cases.Title(language.English).String(user)
	posts := []models.Post{}
	ctx := context.Background()

	// Codes that pulls user from database.
	cur, err := uc.Db.SearchUser(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprintln(w, "Username not registered!")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = cur.All(ctx, &posts)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "An error occured, try again...", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "no record found for user!!")
		return
	}

	json.NewEncoder(w).Encode(posts)
}

// Profile returns the user's profile.
func (uc UserController) Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	post := models.Post{}
	ctx := context.Background()
	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Unable to retrieve user information!")
		return
	}

	find, err := uc.Db.GetUserPosts(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprintln(w, "An error occured, try again... ")
		return
	}

	// Should the post be a slice of post or just a pos?t
	for find.Next(ctx) {
		err := find.Decode(&post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print("error decoding post into struct: ", err)
			fmt.Fprintln(w, "An error occurred, try again...")
			return
		}
		json.NewEncoder(w).Encode(post)
	}

}

// Feed returns a series of post from multiple users.
func (uc UserController) Feed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	w.Header().Set("Content-Type", "application/json")
	post := []models.Post{}
	ctx := context.Background()
	// Code that displays posts from all users
	cur, err := uc.Db.GetPosts()
	if err != nil {
		http.Error(w, "Error getting posts from database", http.StatusInternalServerError)
		return
	}

	err = cur.All(ctx, &post)
	if err != nil {
		http.Error(w, "Error getting posts from database", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)

}

// Logout lets the user log out, deleting the JWT token in the process
func (uc UserController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	// Delete current session and redirect to login page.
	c := &http.Cookie{
		Name:   "",
		Value:  "",
		MaxAge: -1,
	}
	w.Header().Set("Authorization", "Bearer")

	http.SetCookie(w, c)
}

// DeleteUser deletes the user information from the database and associated posts as well.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authetication is still valid
	w.Header().Set("Content-Type", "application/json")
	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "could not retrieve user info, please try again...")
		return
	}
	delete, err := uc.Db.DeleteUser(user)
	if err != nil {
		fmt.Fprintln(w, "error deleting user from collection: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(delete)

	del, err := uc.Db.DeletePosts(user)
	if err != nil {
		fmt.Fprintln(w, "error deleting user's posts from collection: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer ")

	c := &http.Cookie{
		Name:   "",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, c)
	log.Println(del)
	report := "User deleted successfully!"
	json.NewEncoder(w).Encode(report)

}

func (uc UserController) ChangePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	p := models.Password{}
	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "could not retrieve user info, please try again...")
		return
	}
	json.NewDecoder(r.Body).Decode(&p)

	if p.New != p.Confirm {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "new passwords do not match!, try again")
		return
	}

	result, err := uc.Db.ChangePassword(user, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "An error occurred, try again!")
		log.Print(err)
		return
	}

	log.Println(result)

	json.NewEncoder(w).Encode("Password change successful.")
}
