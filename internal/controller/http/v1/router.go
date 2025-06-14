package v1

import (
	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/usecase"
	"github.com/AlexandrKudryavtsev/GoMovieSearch/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, logger logger.Interface, movieUseCase usecase.Movies) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	router := handler.Group("/api")

	newCommonRoutes(router)
	newMoviesRoutes(router, logger, movieUseCase)
}
