package controller

// How to add likes and dislike?
// Do you store in a database? If so, how?
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/A-Victory/blog-API/auth"
	"github.com/julienschmidt/httprouter"
)

// Add the post id to the user database
// For each user, add the column upvote and downvote to each user

func (uc UserController) Upvote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the value from the request
	w.Header().Set("Content-Type", "application/json")
	post_id := ps.ByName("id")

	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error getting user info")
		return
	}

	err = uc.Db.UpVote(post_id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprintln(w, "unable to implement UpVote!")
		return
	}
	// How do I get the user?
	// After getting the user, use it for the filter
	json.NewEncoder(w).Encode("Upvote successfull")
	// Upload to the database the the post id to the appropriate vote column
}

func (uc UserController) Downvote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the value from the request
	w.Header().Set("Content-Type", "application/json")
	post_id := ps.ByName("id")

	user, err := auth.GetUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error getting user info")
		return
	}

	if err := uc.Db.DownVote(post_id, user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		fmt.Fprintln(w, "unable to implement DownVote!")
		return
	}
	// How do I get the user?
	// After getting the user, use it for the filter
	json.NewEncoder(w).Encode("Downvote successfull")
	// Upload to the database the the post id to the appropriate vote column
}
