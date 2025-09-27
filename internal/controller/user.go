package controller

import (
	"go-api/internal/entity"
	"go-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	service service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface) *userController {
	return &userController{service: service}
}

func (c *userController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAll())
}

func (c *userController) Login(ctx *gin.Context) {

	var body entity.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginData, err := c.service.Login(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//here we set the cookie
	ctx.SetCookie(
		"session_id",      // name
		loginData.Session, // value (the session ID from your service)
		3600,              // maxAge in seconds (1 hour)
		"/",               // path
		"localhost",       // domain (use your domain in production)
		false,             // secure (only send over HTTPS)
		true,              // httpOnly (not accessible via JS)
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successfull",
		"data":    loginData.User,
	})
}

func (c *userController) Logout(ctx *gin.Context) {
	session, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You are already logged out",
		})
		return
	}

	err = c.service.Logout(session)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.SetCookie(
		"session_id", // name must match the one you set at login
		"",           // empty value
		-1,           // MaxAge < 0 means delete immediately
		"/",          // path
		"localhost",  // domain
		false,        // secure (HTTPS only in prod)
		true,         // httpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successfull",
	})
}

func (c *userController) RegisterUser(ctx *gin.Context) {
	var body entity.RegisterRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.Register(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": body.Email})
}
