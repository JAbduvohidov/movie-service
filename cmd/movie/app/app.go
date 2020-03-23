package app

import (
	"errors"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	jwt2 "github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/rest"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"movie-service/pkg/core/movie"
	"movie-service/pkg/core/token"
	"net/http"
	"strings"
)

type Server struct {
	router   *mux.ExactMux
	pool     *pgxpool.Pool
	secret   jwt.Secret
	movieSvc *movie.Service
}

func NewServer(router *mux.ExactMux, pool *pgxpool.Pool, secret jwt.Secret, userSvc *movie.Service) *Server {
	return &Server{router: router, pool: pool, secret: secret, movieSvc: userSvc}
}

func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleGetMovies() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		movies, err := s.movieSvc.GetAllMovies()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
		}

		if len(movies) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		err = rest.WriteJSONBody(writer, &movies)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleNewMovie() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dto movie.ResponseDTO
		err := rest.ReadJSONBody(request, &dto)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}

		empty, reason := checkIfEmpty(dto)
		if empty {
			err = rest.WriteJSONBody(writer, reason)
			log.Print(err)
			return
		}

		movi := movie.ResponseDTO{
			Id:          dto.Id,
			Title:       dto.Title,
			Description: dto.Description,
			Image:       dto.Image,
			Year:        dto.Year,
			Country:     dto.Country,
			Actors:      dto.Actors,
			Genres:      dto.Genres,
			Creators:    dto.Creators,
			Studio:      dto.Studio,
			ExtLink:     dto.ExtLink,
		}

		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)
		err = s.movieSvc.AddMovie(payload, request.Context(), movi)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			if errors.Is(err, movie.ErrPrivilege) {
				_ = rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.not_enough_privilege"},
				})
				return
			}
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.bad_request"},
			})
			log.Print(err)
			return
		}
	}
}

func checkIfEmpty(dto movie.ResponseDTO) (ok bool, err ErrorDTO) {
	if strings.TrimSpace(dto.Title) == "" && len(strings.TrimSpace(dto.Title)) == 0 {
		err.Errors = append(err.Errors, "err.invalid_title")
		return true, err
	}

	if strings.TrimSpace(dto.Description) == "" && len(strings.TrimSpace(dto.Description)) == 0 {
		err.Errors = append(err.Errors, "err.invalid_description")
		return true, err
	}

	if strings.TrimSpace(dto.Image) == "" && len(strings.TrimSpace(dto.Image)) == 0 {
		err.Errors = append(err.Errors, "err.invalid_image")
		return true, err
	}
	return false, err
}

func (s *Server) handleGetMovie() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mov, err := s.movieSvc.GetMovie(request)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
		}

		err = rest.WriteJSONBody(writer, &mov)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleSearchMovie() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		searchValue := request.FormValue("q")
		searchValue = strings.Replace(searchValue, "+", "%", -1)
		log.Print(searchValue)
		mov, err := s.movieSvc.SearchMovie(searchValue)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
		}

		if len(mov) == 0 {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		err = rest.WriteJSONBody(writer, &mov)
		if err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) handleDeleteMovie() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		payload := request.Context().Value(jwt2.ContextKey("jwt")).(*token.Payload)

		err := s.movieSvc.RemoveMovie(payload, request)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				writer.WriteHeader(http.StatusNoContent)
				return
			}
			if errors.Is(err, movie.ErrPrivilege) {
				writer.WriteHeader(http.StatusBadRequest)
				_ = rest.WriteJSONBody(writer, &ErrorDTO{
					[]string{"err.not_enough_privilege"},
				})
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			_ = rest.WriteJSONBody(writer, &ErrorDTO{
				[]string{"err.internal_server_error"},
			})
			log.Print(err)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}