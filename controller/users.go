package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/A-Victory/blog-API/auth"
	"github.com/A-Victory/blog-API/db"
	"github.com/A-Victory/blog-API/models"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

	_, _, err = uc.Db.GetUser(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			insert, err := uc.Db.CreateUser(user)
			if err != nil {
				http.Error(w, "Error creating accounting, please try again...", http.StatusInternalServerError)
				return
			}
			fmt.Fprintln(w, insert)

		}
	} else {
		fmt.Fprintln(w, "Email already registered to an account! Please try again...")
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}
	// Send verification token to user email for verification.
	// If no error is encountered, return success.
	fmt.Fprintln(w, "Account successfully created!")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Login allows user to login in to existing account, creating a JWT token in the process.
func (uc UserController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.Header.Set("Content-Type", "application/Json")
	user := models.User{}
	//ctx := context.Background()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	// Check the validity of information provided by the user in database
	if user.Email == "" {
		fmt.Fprintln(w, "Please provide an email address")
	}

	username, passwrd, err := uc.Db.GetUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwrd), []byte(user.Password)); err != nil {
		fmt.Fprintln(w, "Incorrect password!!")
		return
	}
	// Send verification token to user email to create a token that stays valid for 15min of inactivity. After which a new token will have to be generated.
	token, err := auth.GenerateJWT(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Token", token)
	w.Header().Set("Content-Type", "application/json")
	// Write the string to the database?
	// If no error is returned, login is successful

	fmt.Fprintln(w, "Login successful!")
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
		http.Error(w, "Couldn't update user info!", http.StatusInternalServerError)
		return
	}
	// Return success if no error is encountered.]

	fmt.Fprintln(w, "user details successfully updated. ")

}

// Search allows a user search for other users. Returning information including the user's posts.
func (uc UserController) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check to see if authentication is still valid

	w.Header().Set("Content-Type", "application/json")
	user := ps.ByName("user")
	// Codes that pulls user from database.
	result, err := uc.Db.SearchUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(result) == 0 {
		fmt.Fprintln(w, "User does not have any post!")
		return
	} else {
		fmt.Fprintf(w, "%v\n", result)
	}
	// Displays the user posts, comments and other information.
	// Returns success

	json.NewEncoder(w).Encode(result)
}

// Profile returns the user's profile.
func (uc UserController) Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	// Code to display users posts, comments and other information.
	// Get post from post databse using the user_id as the filter for search.
	post := []models.Post{}
	ctx := context.Background()
	user, err := auth.GetUser(r)
	if err != nil {
		fmt.Fprintln(w, "Unable to retrieve user information!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	find, err := uc.Db.GetUserPosts(user)
	if err != nil {
		fmt.Fprintln(w, "cowuld not find user's posts: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Should the post be a slice of post or just a pos?t
	for find.Next(ctx) {
		err := find.Decode(post)
		if err != nil {
			fmt.Fprintln(w, "error decoding post into struct: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(post)
	}

	// Get the user info from the authentication(jwt)
	// Returns success

	fmt.Println(w, "models.User profile is displayed")
}

// Feed returns a series of post from multiple users.
func (uc UserController) Feed(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	post := []models.Post{}
	ctx := context.Background()
	// Code that displays posts from all users
	cur, err := uc.Db.GetPosts()
	if err != nil {
		http.Error(w, "Error getting posts from database", http.StatusInternalServerError)
		return
	}

	for cur.Next(ctx) {
		err := cur.Decode(post)
		if err != nil {
			http.Error(w, "Error getting posts from database", http.StatusInternalServerError)
			return
		}
	}
	// Post are shown, user the empty filter to generate the posts.
	fmt.Fprintf(w, "%v\n", post)

	// Returns success

	fmt.Fprintln(w, "Feed is displayed")
}

// Logout lets the user log out, deleting the JWT token in the process
func (uc UserController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	// Delete current session and redirect to login page.

}

// DeleteUser deletes the user information from the database and associated posts as well.
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authetication is still valid
	user, err := auth.GetUser(r)
	if err != nil {
		fmt.Fprintln(w, "Unable to retrieve user information!")
		w.WriteHeader(http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(del)

	//json.NewEncoder(w).Encode(delete)
	// Should post and comments be deleted as well?, What about users votes?
	// Return success.

	fmt.Fprintln(w, "User deleted successfully!")
}
