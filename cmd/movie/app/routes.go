package app

import (
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
		logger.Logger("MOVIES"),
	)

	s.router.DELETE("/api/movies/{id}",
		s.handleDeleteMovie(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("MOVIES"),
	)

	s.router.POST("/api/movies",
		s.handleNewMovie(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, false, ""),
		jwt.JWT(jwt.SourceAuthorization, reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("MOVIES"),
	)

	s.router.GET("/api/movies/search",
		s.handleSearchMovie(),
		logger.Logger("MOVIES"),
	)
}
