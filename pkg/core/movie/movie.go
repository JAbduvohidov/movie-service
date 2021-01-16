package movie

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"movie-service/pkg/core/token"
	"net/http"
)

var ErrPrivilege = errors.New("not enough privilege")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type ResponseDTO struct {
	Id          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Year        string   `json:"year"`
	Country     string   `json:"country"`
	Actors      []string `json:"actors"`
	Genres      []string `json:"genres"`
	Creators    []string `json:"creators"`
	Studio      string   `json:"studio"`
	ExtLink     string   `json:"ext_link"`
}

type RequestDTO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Year        string   `json:"year"`
	Country     string   `json:"country"`
	Actors      []string `json:"actors"`
	Genres      []string `json:"genres"`
	Creators    []string `json:"creators"`
	Studio      string   `json:"studio"`
	ExtLink     string   `json:"ext_link"`
}

func (s *Service) AddMovie(payload *token.Payload, ctx context.Context, request ResponseDTO) (err error) {
	if payload.Role != "ADMIN" && payload.Role != "MODERATOR" {
		return ErrPrivilege
	}

	if request.Id == 0 {
		_, err = s.pool.Exec(ctx, addMovieDML,
			request.Title,
			request.Description,
			request.Image,
			request.Year,
			request.Country,
			request.Actors,
			request.Genres,
			request.Creators,
			request.Studio,
			request.ExtLink,
		)
	} else {
		_, err = s.pool.Exec(ctx, updateMovieDML,
			request.Title,
			request.Description,
			request.Year,
			request.Country,
			request.Actors,
			request.Genres,
			request.Creators,
			request.Studio,
			request.ExtLink,
			request.Id,
		)
		if err != nil {
			return fmt.Errorf("unable to add movie: %w", err)
		}

		if request.Image != "" {
			_, err = s.pool.Exec(ctx, updateMovieImageDML,
				request.Image,
				request.Id,
			)
		}
	}

	if err != nil {
		return fmt.Errorf("unable to add movie: %w", err)
	}
	return nil
}

func (s *Service) GetAllMovies() (movies []*ResponseDTO, err error) {
	rows, err := s.pool.Query(context.Background(), getAllMoviesDML)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		movi := ResponseDTO{}
		err := rows.Scan(
			&movi.Id,
			&movi.Title,
			&movi.Description,
			&movi.Image,
			&movi.Year,
			&movi.Country,
			&movi.Actors,
			&movi.Genres,
			&movi.Creators,
			&movi.Studio,
			&movi.ExtLink,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movi)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, err
}

func (s *Service) GetMovie(request *http.Request) (movie ResponseDTO, err error) {
	row := s.pool.QueryRow(context.Background(), getMovieDML, request.Context().Value("id"))

	err = row.Scan(
		&movie.Id,
		&movie.Title,
		&movie.Description,
		&movie.Image,
		&movie.Year,
		&movie.Country,
		&movie.Actors,
		&movie.Genres,
		&movie.Creators,
		&movie.Studio,
		&movie.ExtLink,
	)
	if err != nil {
		return
	}

	return movie, nil
}

func (s *Service) SearchMovie(mov string) (movies []ResponseDTO, err error) {
	rows, err := s.pool.Query(context.Background(), searchMoviesDML,
		"%"+mov+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		movi := ResponseDTO{}
		err := rows.Scan(
			&movi.Id,
			&movi.Title,
			&movi.Description,
			&movi.Image,
			&movi.Year,
			&movi.Country,
			&movi.Actors,
			&movi.Genres,
			&movi.Creators,
			&movi.Studio,
			&movi.ExtLink,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, movi)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, err
}

func (s *Service) RemoveMovie(payload *token.Payload, request *http.Request) (err error) {
	if payload.Role != "ADMIN" && payload.Role != "MODERATOR" {
		return ErrPrivilege
	}
	_, err = s.pool.Exec(context.Background(), deleteMovieDML, request.Context().Value("id"))
	if err != nil {
		return err
	}
	return nil
}
