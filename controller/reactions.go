package controller

// How to add likes and dislike?
// Do you store in a database? If so, how?
import (
	"fmt"
	"net/http"

	"github.com/A-Victory/blog-API/auth"
	"github.com/julienschmidt/httprouter"
)

// Add the post id to the user database
// For each user, add the column upvote and downvote to each user

func (uc UserController) Upvote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the value from the request
	post_id := ps.ByName("id")

	user, err := auth.GetUser(r)
	if err != nil {
		fmt.Fprintln(w, "Error getting user info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := uc.Db.UpVote(post_id, user); err != nil {
		fmt.Fprintln(w, "unable to implement DownVote!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// How do I get the user?
	// After getting the user, use it for the filter

	// Upload to the database the the post id to the appropriate vote column
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (uc UserController) Downvote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the value from the request
	post_id := ps.ByName("id")

	user, err := auth.GetUser(r)
	if err != nil {
		fmt.Fprintln(w, "Error getting user info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := uc.Db.DownVote(post_id, user); err != nil {
		fmt.Fprintln(w, "unable to implement DownVote!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// How do I get the user?
	// After getting the user, use it for the filter

	// Upload to the database the the post id to the appropriate vote column
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
