// Package server represents mostly used incoming port.
// There is http server with router, but it could be a grpc server, or tcp with custom protocol.
package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"microservice/internal/app/exibillia"
	errt "microservice/internal/errors"
	"net/http"
)

type ExibilliaService interface {
	Create(ctx context.Context, exibillia exibillia.Exibillia) error
	GetByID(ctx context.Context, id uint64) (exibillia.Exibillia, error)
	Update(ctx context.Context, req exibillia.UpdateRequest) (exibillia.Exibillia, error)
	Delete(ctx context.Context, id uint64) error
}

type HTTPServer struct {
	exibillia ExibilliaService
}

func NewServer(exibillia ExibilliaService) *HTTPServer {
	return &HTTPServer{exibillia: exibillia}
}

func (s *HTTPServer) Listen() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/v1", s.apiRouter())

	return http.ListenAndServe(":3000", r)
}

func (s *HTTPServer) apiRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/exibillia/{exibilliaID}", s.GetExibillia)
	r.Post("/exibillia", s.CreateExibillia)
	r.Put("/exibillia", s.UpdateExibillia)
	r.Delete("/exibillia/{exibilliaID}", s.DeleteExibillia)

	return r
}

func (s *HTTPServer) successResponse(res interface{}, w http.ResponseWriter) {
	if res != nil {
		data, err := json.Marshal(res)

		if err != nil {
			zap.S().Errorw("marshal response", "error", err)

			s.internalError(w)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			zap.S().Errorw("write http response body", "error", err)
		}
	}
}

func (s *HTTPServer) errorResponse(err error, w http.ResponseWriter) {
	switch err.(type) {
	case errt.NotFoundError:
		s.notFoundError(err.Error(), w)
	default:
		s.internalError(w)
	}
}

func (s *HTTPServer) badRequestError(what string, w http.ResponseWriter) {
	res := apiError{Error: what}
	data, err := json.Marshal(res)

	if err != nil {
		zap.S().Errorw("marshal error response", "error", err)

		s.internalError(w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	_, err = w.Write(data)
	if err != nil {
		zap.S().Errorw("write http response body", "error", err)
	}
}

func (s *HTTPServer) notFoundError(what string, w http.ResponseWriter) {
	res := apiError{Error: what}
	data, err := json.Marshal(res)

	if err != nil {
		zap.S().Errorw("marshal error response", "error", err)

		s.internalError(w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, err = w.Write(data)
	if err != nil {
		zap.S().Errorw("write http response body", "error", err)
	}
}

func (s *HTTPServer) internalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *HTTPServer) closeBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		zap.S().Errorw("close http request body", "error", err)
	}
}
