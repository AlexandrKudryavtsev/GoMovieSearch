package v1

import (
	"net/http"

	"github.com/AlexandrKudryavtsev/GoMovieSearch/pkg/logger"
	"github.com/gin-gonic/gin"
)

type moviesRoutes struct {
	l logger.Interface
}

func newMoviesRoutes(handler *gin.RouterGroup, logger logger.Interface) {
	r := moviesRoutes{logger}

	handler.GET("/movie/autocomplete", r.doMovieAutoComplete)
	handler.POST("/movie/index", r.doMovieIndex)
}

// @Summary     Get movie autocomplete
// @Description Get movie autocomplete
// @ID          get-movie-autocomplete
// @Tags  	    movie
// @Accept      json
// @Success     200
// @Failure     500
// @Produce     json
// @Router      /movie/autocomplete [get]
func (r *moviesRoutes) doMovieAutoComplete(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

// func
// @Summary     Index movies
// @Description Index movies
// @ID          post-movie-index
// @Tags  	    movie
// @Accept      json
// @Success     200
// @Failure     500
// @Produce     json
// @Router      /movie/index [post]
func (r *moviesRoutes) doMovieIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
