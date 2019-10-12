package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jogeocc/go-mysql-crud/driver"
	models "github.com/jogeocc/go-mysql-crud/models"
	repository "github.com/jogeocc/go-mysql-crud/repository"
	pelicula "github.com/jogeocc/go-mysql-crud/repository/pelicula"
)

// NewPeliculaHandler ...
func NewPeliculaHandler(db *driver.DB) *Pelicula {
	return &Pelicula{
		repo: pelicula.NewSQLPeliculaRepo(db.SQL),
	}
}

// Estrucura nueva basada en el modelo...
type Pelicula struct {
	repo repository.PeliculaMethods
}

// Fetch all post data
func (p *Pelicula) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, _ := p.repo.Fetch(r.Context(), 5)

	respondwithJSON(w, http.StatusOK, payload)
}

// Create a new post
func (p *Pelicula) Create(w http.ResponseWriter, r *http.Request) {
	pelicula := models.Pelicula{}
	json.NewDecoder(r.Body).Decode(&pelicula)

	newID, err := p.repo.Create(r.Context(), &pelicula)
	fmt.Println(newID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

// Update a post by id
func (p *Pelicula) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	data := models.Pelicula{ID: int(id)}
	json.NewDecoder(r.Body).Decode(&data)
	payload, err := p.repo.Update(r.Context(), &data)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// GetByID returns a post details
func (p *Pelicula) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := p.repo.GetByID(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Content not found")
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// Delete a post
func (p *Pelicula) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := p.repo.Delete(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusMovedPermanently, map[string]string{"message": "Delete Successfully"})
}

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}
