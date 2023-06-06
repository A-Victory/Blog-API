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

	r.GET("/search/:user", auth.Verify(uc.Search))
	r.GET("/profile", auth.Verify(uc.Profile))
	r.PATCH("/user", auth.Verify(uc.UpdateInfo))
	r.GET("/logout", auth.Verify(uc.Logout))
	r.GET("/feed", auth.Verify(uc.Feed))
	r.DELETE("/user", auth.Verify(uc.DeleteUser))

	r.POST("/post", auth.Verify(uc.CreatePost))
	r.DELETE("/post/:id", auth.Verify(uc.DeletePost))
	r.PUT("/post/:id", auth.Verify(uc.EditPost))
	r.GET("/post/:id/", auth.Verify(uc.ViewPost))
	r.DELETE("/post/comment/:id", auth.Verify(uc.DeleteComment))
	r.POST("/post/comment/:id", auth.Verify(uc.AddComment))

	r.GET("/upvote/:id", auth.Verify(uc.Upvote))
	r.GET("/downvote/:id", auth.Verify(uc.Downvote))

}
