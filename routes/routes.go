package routes

import (
	"github.com/A-Victory/blog-API/auth"
	"github.com/A-Victory/blog-API/controller"
	"github.com/A-Victory/blog-API/db"
	"github.com/julienschmidt/httprouter"
)

var Routers = func(r *httprouter.Router) {
	uc := controller.NewUserController(db.UserDb())
	r.POST("/signup", uc.Signup)
	r.POST("/login", uc.Login)
	r.GET("/profile", auth.Verify(uc.Profile))
	r.PATCH("/user", auth.Verify(uc.UpdateInfo))
	r.GET("/logout", auth.Verify(uc.Logout))

}
