package v1

import (
	"net/http"

	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/entity"
	"github.com/AlexandrKudryavtsev/GoMovieSearch/internal/usecase"
	"github.com/AlexandrKudryavtsev/GoMovieSearch/pkg/logger"
	"github.com/gin-gonic/gin"
)

type moviesRoutes struct {
	u usecase.Movies
	l logger.Interface
}

func newMoviesRoutes(handler *gin.RouterGroup, logger logger.Interface, moviesUseCase usecase.Movies) {
	r := moviesRoutes{moviesUseCase, logger}

	handler.POST("/movies/index", r.doIndexMovies)
	handler.GET("/movies/search", r.doSearchMovies)
	handler.GET("/movies/autocomplete", r.doAutocompleteMovies)
}

type doIndexMoviesRequest struct {
	Movies []entity.Movie `json:"movies" binding:"required"`
}

// @Summary     Index movies
// @Description Add or update movies in search index
// @ID          index-movies
// @Tags  	    movies
// @Param       request body doIndexMoviesRequest true "Movies data"
// @Accept      json
// @Success     200
// @Failure     400
// @Failure     500
// @Produce     json
// @Router      /movies/index [post]
func (m *moviesRoutes) doIndexMovies(ctx *gin.Context) {
	var request doIndexMoviesRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		m.l.Error(err, "http - v1 - doIndexMovies")
		errorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := m.u.Index(ctx.Request.Context(), request.Movies); err != nil {
		m.l.Error(err, "http - v1 - doIndexMovies")
		errorResponse(ctx, http.StatusInternalServerError, "failed to index movies")
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"count":  len(request.Movies),
	})
}

// @Summary     Search movies
// @Description Search movies by query
// @ID          search-movies
// @Tags  	    movies
// @Param       query query string true "Search query"
// @Accept      json
// @Success     200 {object} map[string]interface{} "data: []entity.Movie"
// @Failure     400
// @Failure     500
// @Produce     json
// @Router      /movies/search [get]
func (m *moviesRoutes) doSearchMovies(ctx *gin.Context) {
	query := ctx.Query("query")
	if len(query) < 2 {
		errorResponse(ctx, http.StatusBadRequest, "query too short")
		return
	}

	movies, err := m.u.Search(ctx.Request.Context(), query)
	if err != nil {
		m.l.Error(err, "http - v1 - doSearchMovies")
		errorResponse(ctx, http.StatusInternalServerError, "search failed")
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": movies,
	})
}

// @Summary     Autocomplete movies
// @Description Get movie suggestions for autocomplete
// @ID          autocomplete-movies
// @Tags  	    movies
// @Param       query query string true "Autocomplete query"
// @Accept      json
// @Success     200 {object} map[string]interface{} "data: []entity.Movie"
// @Failure     400
// @Failure     500
// @Produce     json
// @Router      /movies/autocomplete [get]
func (m *moviesRoutes) doAutocompleteMovies(ctx *gin.Context) {
	query := ctx.Query("query")
	if len(query) < 2 {
		errorResponse(ctx, http.StatusBadRequest, "query too short")
		return
	}

	movies, err := m.u.Autocomplete(ctx.Request.Context(), query)
	if err != nil {
		m.l.Error(err, "http - v1 - doAutocompleteMovies")
		errorResponse(ctx, http.StatusInternalServerError, "autocomplete failed")
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": movies,
	})
}
