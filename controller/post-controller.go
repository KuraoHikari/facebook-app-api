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

type PostController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type postController struct {
	postService service.PostService
	jwtService  service.JWTService
}

func NewPostController(postServ service.PostService, jwtServ service.JWTService) PostController {
	return &postController{
		postService: postServ,
		jwtService:  jwtServ,
	}
}

func (c *postController) All(context *gin.Context){
	var posts []entity.Post = c.postService.All()
	res := helper.BuildSuccessResponse(true, "OK", posts)
	context.JSON(http.StatusOK, res)
}

func (c *postController) FindByID(context *gin.Context){
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var post entity.Post = c.postService.FindByID(id)
	if (post == entity.Post{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildSuccessResponse(true, "OK", post)
		context.JSON(http.StatusOK, res)
	}
}
func (c *postController) Insert(context *gin.Context){
	var postCreateDTO dto.PostCreateDTO
	errDTO := context.ShouldBind(&postCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		// userID := c.getUserIDByToken(authHeader)
		aToken, err := c.jwtService.ValidateToken(authHeader)
			if err != nil {
				res := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
				context.AbortWithStatusJSON(http.StatusBadRequest, res)
				return
			}
		claims := aToken.Claims.(jwt.MapClaims)
		userID := fmt.Sprintf("%v", claims["user_id"])
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			postCreateDTO.UserID = convertedUserID
		}
		result := c.postService.Insert(postCreateDTO)
		response := helper.BuildSuccessResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}
func (c *postController) Update(context *gin.Context){
	var postUpdateDTO dto.PostUpdateDTO
	errDTO := context.ShouldBind(&postUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
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
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.postService.IsAllowedToEdit(userID, postUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			postUpdateDTO.UserID = id
		}
		result := c.postService.Update(postUpdateDTO)
		response := helper.BuildSuccessResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}
func (c *postController) Delete(context *gin.Context){
	var post entity.Post
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	post.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		res := helper.BuildErrorResponse("Failed to process request", errToken.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.postService.IsAllowedToEdit(userID, post.ID) {
		c.postService.Delete(post)
		res := helper.BuildSuccessResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}