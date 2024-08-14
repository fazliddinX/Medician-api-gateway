package user

import (
	"api-gateway/generated/users"
	"api-gateway/pkg/models"
	"api-gateway/pkg/token"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"

	_ "api-gateway/api/docs"
)

type UserImpl struct {
	User users.UserServiceClient
	Log  *slog.Logger
}

// @Summary Get User Profile
// @Description Retrieves the user profile based on JWT claims
// @Tags User
// @Produce json
// @Success 200 {object} users.UserResponse
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/user [get]
func (u *UserImpl) GetUserProfile(c *gin.Context) {

	value, check := c.Get("claims")
	if !check {
		u.Log.Error("Claims not found")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "Claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		u.Log.Error("cannot extrackt to clams")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "Claims not extracted"})
		return
	}

	req := users.UserID{Id: claims.ID}

	resp, err := u.User.GetUserProfile(context.Background(), &req)
	if err != nil {
		u.Log.Error("GetUserProfile error", "error", err)
		c.JSON(http.StatusUnauthorized, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Update User
// @Description Updates user information based on the provided request body and JWT claims
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUserRequest body users.UpdateUserRequest true "Update User Request"
// @Success 200 {object} users.UserResponse
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/user [put]
func (u *UserImpl) UpdateUser(c *gin.Context) {

	value, check := c.Get("claims")
	if !check {
		u.Log.Error("Claims not found")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		u.Log.Error("cannot extrackt to clams")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "Claims not extracted"})
		return
	}

	var req users.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		u.Log.Error("ShouldBindJSON error", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.Id = claims.ID

	res, err := u.User.UpdateUser(context.Background(), &req)
	if err != nil {
		u.Log.Error("UpdateUser error", "error", err)
		c.JSON(http.StatusUnauthorized, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Delete User
// @Description Deletes the user based on JWT claims
// @Tags User
// @Produce json
// @Success 200 {object} users.Message
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/user [delete]
func (u *UserImpl) DeleteUser(c *gin.Context) {

	value, check := c.Get("claims")
	if !check {
		u.Log.Error("Claims not found")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "Claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		u.Log.Error("cannot extrackt to clams")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "Claims not extracted"})
		return
	}

	req := users.UserID{Id: claims.ID}

	res, err := u.User.DeleteUser(context.Background(), &req)
	if err != nil {
		u.Log.Error("DeleteUser error", "error", err)
		c.JSON(http.StatusUnauthorized, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Get All Users
// @Description Retrieves a list of users based on query parameters
// @Tags User
// @Produce json
// @Param limit query string false "Limit the number of users returned"
// @Param offset query string false "Offset for pagination"
// @Param firstName query string false "Filter by first name"
// @Param gender query string false "Filter by gender"
// @Success 200 {object} users.GetAllUsersResponse
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/user/all [get]
func (u *UserImpl) GetAllUsers(c *gin.Context) {
	var filter models.Filter

	if err := c.ShouldBindQuery(&filter); err != nil {
		u.Log.Error("ShouldBindQuery error", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	var limit int
	var offset int
	var err error
	if filter.Limit != "" {
		limit, err = strconv.Atoi(filter.Limit)
		if err != nil {
			u.Log.Error("strconv error limit", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}
	}

	if filter.Offset != "" {
		offset, err = strconv.Atoi(filter.Offset)
		if err != nil {
			u.Log.Error("strconv error offset", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}
	}

	req := users.GetAllUsersRequest{
		Limit:     int64(limit),
		Offset:    int64(offset),
		FirstName: filter.FirstName,
		Gender:    filter.Gender,
	}

	resp, err := u.User.GetAllUsers(context.Background(), &req)
	if err != nil {
		u.Log.Error("GetAllUsers error", "error", err)
		c.JSON(http.StatusUnauthorized, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Create User
// @Description Creates a new user with the provided details
// @Tags User
// @Accept json
// @Produce json
// @Param CreateUserRequest body users.CreateUserRequest true "Create User Request"
// @Success 200 {object} users.CreateUserResponse
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/user/create [post]
func (u *UserImpl) CreateUser(c *gin.Context) {
	var req users.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		u.Log.Error("ShouldBindJSON error", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	res, err := u.User.CreateUser(context.Background(), &req)
	if err != nil {
		u.Log.Error("CreateUser error", "error", err)
		c.JSON(http.StatusUnauthorized, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
