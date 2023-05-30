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
	// Write code that saves the user information in database.
	/*
		coll := uc.Db.Collection("users")
		insert, err := coll.InsertOne(context.Background(), user)
		if err != nil {
			fmt.Fprintln(w, "Unable to insert user")
			return
		}
	*/
	err = uc.Va.ValidateUserInfo(user)
	if err != nil {
		http.Error(w, "Error validating user infomation.....", http.StatusInternalServerError)
		return
	}

	insert, err := uc.Db.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, insert)
	// Send verification token to user email for verification.
	// If no error is encountered, return success.
	fmt.Fprintln(w, "Account created!!")
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
	/*
		filter := bson.D{{Key: "email", Value: user.Email}}
		coll := uc.Db.Collection("users")
		find := coll.FindOne(ctx, filter)

		info := models.User{}
		err := find.Decode(&info)
		if err != nil {
			fmt.Fprintln(w, "Account does not exist, please signup for a new account...")
			return
		}
	*/
	passwrd, err := uc.Db.GetUser(user)
	if err == db.ErrDb {
		http.Error(w, "Error getting user information!", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwrd), []byte(user.Password)); err != nil {
		fmt.Fprintln(w, "Incorrect password!!")
		return
	}
	// Send verification token to user email to create a token that stays valid for 15min of inactivity. After which a new token will have to be generated.
	token, err := auth.GenerateJWT(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Token", token)
	w.Header().Set("Content-Type", "application/json")
	// Write the string to the database?
	// If no error is returned, login is successful

	fmt.Fprintln(w, "Login successful")
}

// UpdateInfo updates the user resources in the database.
func (uc UserController) UpdateInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check to see if authentication is still valid
	r.Header.Set("Content-Type", "application/Json")
	user := models.User{}
	//ctx := context.Background()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	// Update fields provided by user.
	/*
		filter := bson.D{{Key: "user", Value: user.Username}}

		if user.Firstname != "" {
			coll := uc.Db.Collection("users")
			update, _ := coll.UpdateOne(ctx, filter, user.Firstname)
			json.NewEncoder(w).Encode(update)
		}
		if user.Lastname != "" {
			coll := uc.Db.Collection("users")
			update, _ := coll.UpdateOne(ctx, filter, user.Lastname)
			json.NewEncoder(w).Encode(update)
		}
		if user.Password != "" {
			coll := uc.Db.Collection("users")
			update, _ := coll.UpdateOne(ctx, filter, user.Password)
			json.NewEncoder(w).Encode(update)
		}

		// In case of email, validation code should be sent to the new email before updating in the database.
		if user.Email != "" {
			coll := uc.Db.Collection("users")
			update, _ := coll.UpdateOne(ctx, filter, user.Firstname)
			json.NewEncoder(w).Encode(update)
		}
	*/
	err := uc.Db.UpdateUser(w, &user)
	if err != nil {
		http.Error(w, "Couldn't update user info!", http.StatusInternalServerError)
		return
	}
	// Return success if no error is encountered.]
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintln(w, "models.User records updated successfully")

}

// Search allows a user search for other users. Returning information including the user's posts.
func (uc UserController) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check to see if authentication is still valid

	user := ps.ByName("user")
	// Codes that pulls user from database.
	result, err := uc.Db.SearchUser(user)
	if err != nil {
		http.Error(w, "Couldn't update user info!", http.StatusInternalServerError)
		return
	}
	if result == nil {
		fmt.Fprintln(w, "User does not have any post!")
	} else {
		fmt.Fprintf(w, "%v\n", result)
	}
	// Displays the user posts, comments and other information.
	// Returns success

	fmt.Fprintln(w, "models.User is displayed")
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

	}
	find, err := uc.Db.GetUserPosts(user)
	if err != nil {

	}

	// Should the post be a slice of post or just a pos?t
	for find.Next(ctx) {
		err := find.Decode(post)
		if err != nil {

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

	}
	delete, err := uc.Db.DeleteUser(user)
	if err != nil {

	}
	json.NewEncoder(w).Encode(delete)

	del, err := uc.Db.DeletePosts(user)
	if err != nil {

	}
	json.NewEncoder(w).Encode(del)

	//json.NewEncoder(w).Encode(delete)
	// Should post and comments be deleted as well?, What about users votes?
	// Return success.

	fmt.Fprintln(w, "models.User deleted successfully")
}
