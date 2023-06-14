package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/A-Victory/blog-API/auth"
	"github.com/A-Victory/blog-API/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreatePost creates a new post
func (uc UserController) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	post := &models.Post{}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Unable to retrieve user information!")
		return
	}

	// Write code that stores the post in database.
	insert, err := uc.Db.CreatePost(user, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(insert.InsertedID)

	json.NewEncoder(w).Encode("Successfully created post, with post ID: " + insert.InsertedID.(primitive.ObjectID).Hex())
}

// DeletePost deletes the post from the database.
func (uc UserController) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("id")

	// Delete Post from database.
	delete, err := uc.Db.DeletePost(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "An error occurred while deleting post, try again...")
		return
	}

	log.Println(delete)

	json.NewEncoder(w).Encode("Successfully deleted post with id: " + id)

}

// EditPost edits an already existing post
func (uc UserController) EditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	post := models.Post{}
	id := ps.ByName("id")

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding JSON")
		return
	}

	err := uc.Db.UpdatePost(w, id, post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "An error occurred while updating post, please try again...")
		return
	}

	json.NewEncoder(w).Encode("Successfuly updated post with id: " + id)

}

// ViewPost retrieves a particular post
func (uc UserController) ViewPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	post := models.Post{}
	id := ps.ByName("id")

	res := uc.Db.GetPost(id)
	if err := res.Decode(&post); err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Id does not match an existing post!")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			fmt.Fprintln(w, "Error encountered while getting post from database, please try again...")
			return
		}
	}

	//Display the post.
	json.NewEncoder(w).Encode(post)
}

// AddComment adds a new comment to a post
func (uc UserController) AddComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("id")
	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Unable to retrieve user information!")
		return
	}
	com := models.Comment{}

	if err := json.NewDecoder(r.Body).Decode(&com); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error while decoding JSON")
		return
	}

	// Attach the comment to the post sending all to the database.
	upload, err := uc.Db.Comment(user, id, com)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprintln(w, "Error while uploading comment to database, please try again...")
		return
	}

	log.Println(upload)

	json.NewEncoder(w).Encode("Comment has been uploaded to post with id: " + id + "!")
}

// DeleteComment deletes comment from post.
func (uc UserController) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	id := ps.ByName("id")
	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Unable to retrieve user information!")
		return
	}

	delete, err := uc.Db.DeleteComment(id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprintln(w, "Error while deleting comment from database, please try again...")
		return
	}

	log.Println(delete)

	json.NewEncoder(w).Encode("Post with id: " + id + " has been successfully deleted!")

}
