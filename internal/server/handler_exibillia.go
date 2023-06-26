package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"microservice/internal/app/exibillia"
	"net/http"
	"strconv"
	"time"
)

type exibilliaAPIModel struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *HTTPServer) GetExibillia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "exibilliaID"), 10, 64)
	if err != nil {
		s.badRequestError("incorrect id", w)

		return
	}
	obj, err := s.exibillia.GetByID(r.Context(), id)
	if err != nil {
		s.errorResponse(err, w)

		return
	}

	res := exibilliaAPIModel{
		ID:          obj.ID,
		Name:        obj.Name,
		Description: obj.Description,
		Tags:        obj.Tags,
		CreatedAt:   obj.CreatedAt,
		UpdatedAt:   obj.UpdatedAt,
	}
	s.successResponse(res, w)
}

func (s *HTTPServer) CreateExibillia(w http.ResponseWriter, r *http.Request) {
	defer s.closeBody(r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.internalError(w)

		return
	}

	var req exibilliaAPIModel
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.badRequestError("invalid json", w)

		return
	}

	model := exibillia.Exibillia{
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
	}

	err = s.exibillia.Create(r.Context(), model)
	if err != nil {
		s.internalError(w)

		return
	}

	s.successResponse(nil, w)
}

func (s *HTTPServer) UpdateExibillia(w http.ResponseWriter, r *http.Request) {
	defer s.closeBody(r)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.errorResponse(err, w)

		return
	}

	var req exibilliaAPIModel
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.badRequestError("invalid json", w)

		return
	}

	updateReq := exibillia.UpdateRequest{
		ID:          req.ID,
		Description: req.Description,
		Tags:        req.Tags,
	}
	updated, err := s.exibillia.Update(r.Context(), updateReq)
	if err != nil {
		s.errorResponse(err, w)

		return
	}

	s.successResponse(updated, w)
}

func (s *HTTPServer) DeleteExibillia(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "exibilliaID"), 10, 64)
	if err != nil {
		s.badRequestError("incorrect id", w)

		return
	}

	err = s.exibillia.Delete(r.Context(), id)
	if err != nil {
		s.errorResponse(err, w)

		return
	}

	s.successResponse(nil, w)
}
