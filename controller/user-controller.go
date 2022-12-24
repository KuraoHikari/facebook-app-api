package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KuraoHikari/facebook-app-res-api/dto"
	"github.com/KuraoHikari/facebook-app-res-api/entity"
	"github.com/KuraoHikari/facebook-app-res-api/helper"
	"github.com/KuraoHikari/facebook-app-res-api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Update(contex *gin.Context)
	Profile(contex *gin.Context)
}
type userController struct {
	userService service.UserService
	jwtService service.JWTService
}
func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService: jwtService,
	}
}

func (c *userController) Login(ctx *gin.Context){
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authResult := c.userService.VerifyCredential(loginDTO.Email,loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID,10))
		v.Token = generatedToken
		res := helper.BuildSuccessResponse(true, "OK!", v)
		ctx.JSON(http.StatusOK, res)
		return
	}
	res :=helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
}
func (c *userController) Register(ctx *gin.Context){
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if !c.userService.IsDuplicateEmail(registerDTO.Email){
		res := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, res)
	}else {
		createdUser := c.userService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		res := helper.BuildSuccessResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusCreated, res)
	}
}
func (c *userController) Update(context *gin.Context) {
	var userUpdateDTO dto.UserUpdateDTO
	errDTO := context.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		res := helper.BuildErrorResponse("Failed to process request", errToken.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	userUpdateDTO.ID = id
	u := c.userService.Update(userUpdateDTO)
	res := helper.BuildSuccessResponse(true, "OK!", u)
	context.JSON(http.StatusOK, res)
}
func (c *userController) Profile(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}else {
		claims := token.Claims.(jwt.MapClaims)
		id := fmt.Sprintf("%v", claims["user_id"])
		user := c.userService.Profile(id)
		res := helper.BuildSuccessResponse(true, "OK", user)
		context.JSON(http.StatusOK, res)
	}
	
}