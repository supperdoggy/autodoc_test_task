package handler

import (
	"context"
	"net/http"

	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/models"
	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IHandler interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	ListUsers(c *gin.Context)
}

type handler struct {
	logger  *zap.Logger
	service service.IService
}

func NewHandler(l *zap.Logger, s service.IService) IHandler {
	return &handler{
		logger:  l,
		service: s,
	}
}

func (h *handler) CreateUser(c *gin.Context) {
	var (
		req  models.CreateUserRequest
		resp models.UserResponse
		ctx  context.Context
	)
	if err := c.Bind(&req); err != nil {
		h.logger.Error("error Bing", zap.Error(err))
		resp.Error = "error reading request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	h.logger.Info("CreateUser", zap.Any("req", req))

	user, err := h.service.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		h.logger.Error("error CreateUser", zap.Error(err), zap.Any("req", req))
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.User = user
	c.JSON(http.StatusOK, resp)
}

func (h *handler) GetUser(c *gin.Context) {
	var (
		req  models.GetUserRequest
		resp models.UserResponse
		ctx  context.Context
		err  error
	)
	if err := c.Bind(&req); err != nil {
		h.logger.Error("error Bing", zap.Error(err))
		resp.Error = "error reading request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	h.logger.Info("GetUser", zap.Any("req", req))

	if req.ID != nil {
		resp.User, err = h.service.GetUserByID(ctx, *req.ID)
	} else if req.Email != nil {
		resp.User, err = h.service.GetUserByID(ctx, *req.Email)
	} else {
		resp.Error = "error invalid request body"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err != nil {
		h.logger.Error("error getting user", zap.Error(err), zap.Any("req", req))
		resp.Error = "error getting user"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) DeleteUser(c *gin.Context) {
	var (
		req  models.DeleteUserRequest
		resp models.ErrorResponse
		ctx  context.Context
		err  error
	)
	if err := c.Bind(&req); err != nil {
		h.logger.Error("error Bing", zap.Error(err))
		resp.Error = "error reading request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	h.logger.Info("DeleteUser", zap.Any("req", req))

	err = h.service.DeleteUser(ctx, req.ID)
	if err != nil {
		h.logger.Error("error deleting user", zap.Error(err), zap.Any("req", req))
		resp.Error = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) ListUsers(c *gin.Context) {
	var (
		req  models.ListUsersRequest
		resp models.ListUserResponse
		ctx  context.Context
		err  error
	)

	if err := c.Bind(&req); err != nil {
		h.logger.Error("error Bing", zap.Error(err))
		resp.Error = "error reading request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	h.logger.Info("ListUsers")

	users, err := h.service.ListUsers(ctx, req.Limit, req.Skip)
	if err != nil {
		h.logger.Error("error listing users", zap.Error(err), zap.Any("req", req))
		resp.Error = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Users = users

	c.JSON(http.StatusOK, resp)
}
