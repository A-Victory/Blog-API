package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/A-Victory/blog-API/auth"
	"github.com/A-Victory/blog-API/models"
	"github.com/julienschmidt/httprouter"
)

// CreatePost creates a new post
func (uc UserController) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	post := &models.Post{}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(post); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := auth.GetUser(r)
	json.NewEncoder(w).Encode(user)
	if err != nil {
		fmt.Fprintln(w, "Unable to retrieve user information!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Code that stores the user_id in the post struct.
	// Code that get the user_id and stores it in the post struct.

	// Write code that stores the post in database.
	insert, err := uc.Db.CreatePost(user, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(insert)
}

// DeletePost deletes the post from the database.
func (uc UserController) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get post from database using post id
	// verify the id
	id := ps.ByName("id")
	//id := r.URL.Query().Get("id")

	// Delete Post from database.
	delete, err := uc.Db.DeletePost(id)
	if err != nil {
		fmt.Fprintln(w, "Failed to delete post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delete)

}

// EditPost edits an already existing post
func (uc UserController) EditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get post from database
	post := models.Post{}
	id := ps.ByName("id")
	//id := r.URL.Query().Get("id")

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	update, err := uc.Db.UpdatePost(id, post)
	if err != nil {
		fmt.Fprintln(w, "Failed to update post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(update)
	// Send the updated fields to the database

}

// ViewPost retrieves a particular post
func (uc UserController) ViewPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get post from database
	post := models.Post{}
	id := ps.ByName("id")
	//id := r.URL.Query().Get("id")

	res := uc.Db.GetPost(id)
	if err := res.Decode(&post); err != nil {
		fmt.Fprintln(w, "Error decoding post from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Display the post.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// AddComment adds a new comment to a post
func (uc UserController) AddComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the post from database
	id := ps.ByName("id")
	//id := r.URL.Query().Get("id")

	com := models.Comment{}

	if err := json.NewDecoder(r.Body).Decode(&com); err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Attach the comment to the post sending all to the database.
	upload, err := uc.Db.Comment(id, com)
	if err != nil {
		fmt.Fprintln(w, "Error uploading comment to database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(upload)
}

// DeleteComment deletes comment from post.
func (uc UserController) DeleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	user, err := auth.GetUser(r)
	if err != nil {
		fmt.Fprintln(w, "Unable to retrieve user information!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	delete, err := uc.Db.DeleteComment(id, user)
	if err != nil {
		fmt.Fprintln(w, "Error deleting comment from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delete)

}
