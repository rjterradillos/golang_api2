package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rjterradillos/golang_api2/dto"
	"github.com/rjterradillos/golang_api2/entity"
	"github.com/rjterradillos/golang_api2/helper"
	"github.com/rjterradillos/golang_api2/service"
)

// CourseController is a ...
type CourseController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type courseController struct {
	courseService service.CourseService
	jwtService    service.JWTService
}

// NewCourseController create a new instances of BoookController
func NewCourseController(courseServ service.CourseService, jwtServ service.JWTService) CourseController {
	return &courseController{
		courseService: courseServ,
		jwtService:    jwtServ,
	}
}

func (c *courseController) All(context *gin.Context) {
	var courses []entity.Course = c.courseService.All()
	res := helper.BuildResponse(true, "OK", courses)
	context.JSON(http.StatusOK, res)
}

func (c *courseController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var course entity.Course = c.courseService.FindByID(id)
	if (course == entity.Course{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", course)
		context.JSON(http.StatusOK, res)
	}
}

func (c *courseController) Insert(context *gin.Context) {
	var courseCreateDTO dto.CourseCreateDTO
	errDTO := context.ShouldBind(&courseCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			courseCreateDTO.UserID = convertedUserID
		}
		result := c.courseService.Insert(courseCreateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *courseController) Update(context *gin.Context) {
	var courseUpdateDTO dto.CourseUpdateDTO
	errDTO := context.ShouldBind(&courseUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.courseService.IsAllowedToEdit(userID, courseUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			courseUpdateDTO.UserID = id
		}
		result := c.courseService.Update(courseUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *courseController) Delete(context *gin.Context) {
	var course entity.Course
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	course.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.courseService.IsAllowedToEdit(userID, course.ID) {
		c.courseService.Delete(course)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *courseController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
