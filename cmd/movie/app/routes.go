package app

import (
	"context"
	"github.com/JAbduvohidov/mux/middleware/authenticated"
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/mux/middleware/logger"
	"movie-service/pkg/core/token"
	"reflect"
)

func (s *Server) InitRoutes() {
	s.router.GET("/api/movies",
		s.handleGetMovies(),
		logger.Logger("MOVIES"),
	)

	s.router.GET("/api/movies/{id}",
		s.handleGetMovie(),
		logger.Logger("MOVIE"),
	)

	s.router.DELETE("/api/movies/{id}",
		s.handleDeleteMovie(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("MOVIE"),
	)

	s.router.POST("/api/movies",
		s.handleNewMovie(),
		authenticated.Authenticated(func(ctx context.Context) bool { return !jwt.IsContextNonEmpty(ctx) }, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("MOVIE"),
	)

	s.router.GET("/api/movies/search",
		s.handleSearchMovie(),
		logger.Logger("SEARCH_MOVIE"),
	)
}
