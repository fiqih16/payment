package controller

import (
	"api-payment/entity"
	"api-payment/helper"
	"api-payment/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	jwtService service.JWTService
	authService service.AuthService
}

func NewAuthController(jwtService service.JWTService, authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
		jwtService: jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var login entity.User
	err := ctx.ShouldBind(&login)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authResponse := c.authService.VeryfiyCredential(login.Email, login.Password)
	if v, ok := authResponse.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.Itoa(int(v.ID)))
		v.Token = generatedToken
		res := helper.BuildResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var user entity.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if c.authService.IsDuplicateEmail(user.Email) {
		res := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	createdUser := c.authService.CreateUser(user)
	token := c.jwtService.GenerateToken(strconv.Itoa(int(createdUser.ID)))
	createdUser.Token = token
	res := helper.BuildResponse(true, "OK!", createdUser)
	ctx.JSON(http.StatusCreated, res)
}
