package msredis

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

// TODO: create some validators
type Server struct {
	log     *slog.Logger
	service *Service
	*gin.Engine
}

func NewServer(log *slog.Logger) *Server {
	server := gin.Default()
	service := NewService(log)

	return &Server{
		log:     log,
		service: service,
		Engine:  server,
	}
}

func (s *Server) RegisterRoutes() {
	s.GET("/user/:id", s.getUserById)
	s.POST("/user", s.createUser)
	s.DELETE("/user/:id", s.deleteUser)
}

func (s *Server) getUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	user, err := s.service.GetUserById(int32(id))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, &user)
}

func (s *Server) createUser(ctx *gin.Context) {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := ctx.ShouldBindJSON(&reqBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	user, err := s.service.CreateUser(reqBody.Email, reqBody.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (s *Server) deleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	user, err := s.service.DeleteUserById(int32(id))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, &user)
}
