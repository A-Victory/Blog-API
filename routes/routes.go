package routes

import (
	"github.com/A-Victory/blog-API/auth"
	"github.com/A-Victory/blog-API/controller"
	"github.com/A-Victory/blog-API/db"
	"github.com/julienschmidt/httprouter"
)

var Routers = func(r *httprouter.Router) {
	db := db.UserDb()
	va := auth.NewValidator()
	uc := controller.NewUserController(db, va)

	r.POST("/signup", uc.Signup)
	r.POST("/login", uc.Login)

	// requires authentication
	r.POST("/changepassword", auth.Verify(uc.ChangePassword))

	r.GET("/search/:name", auth.Verify(uc.Search))
	r.GET("/profile", auth.Verify(uc.Profile))
	r.PATCH("/user", auth.Verify(uc.UpdateInfo))
	r.GET("/logout", auth.Verify(uc.Logout))
	r.GET("/feed", auth.Verify(uc.Feed))
	r.DELETE("/user", auth.Verify(uc.DeleteUser))
	r.POST("/password", auth.Verify(uc.ChangePassword))

	r.POST("/post", auth.Verify(uc.CreatePost))
	r.DELETE("/post/:id", auth.Verify(uc.DeletePost))
	r.PUT("/post/:id", auth.Verify(uc.EditPost))
	r.GET("/post/:id", auth.Verify(uc.ViewPost))

	r.DELETE("/comment/:id", auth.Verify(uc.DeleteComment))
	r.POST("/comment/:id", auth.Verify(uc.AddComment))

	r.GET("/upvote/:id", auth.Verify(uc.Upvote))
	r.GET("/downvote/:id", auth.Verify(uc.Downvote))

}
