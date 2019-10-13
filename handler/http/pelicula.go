package handler

/*
CLASE QUE PDEFINE LOS METODOS HTTP DONDE RECIBIRA LOS METODOS DESDE LA REQUEST o PETICION WEB Y
LLAMARA LOS METODOS DE LA BD PARA EFECTUAR LA ACCION
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	pelicula "../../repository/pelicula"

	repository "../../repository"

	models "../../models"

	"../../driver"
	"github.com/go-chi/chi"
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

// TRAER TODOS LAS PELICULAS
func (p *Pelicula) Fetch(w http.ResponseWriter, r *http.Request) {
	payload, err := p.repo.Fetch(r.Context())

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusOK, payload)
}

func (p *Pelicula) Create(w http.ResponseWriter, r *http.Request) { // CREAR UNA NUEVA PELI
	pelicula := models.Pelicula{} //CREAMOS UN NEW OBJ DESDE EL MAP DEL MODELO

	r.ParseForm()                                        // CONVERTIMOS LA REQUEST A UN FORMULARIO
	pelicula.Nombre = r.FormValue("nombre")              //INSERTAMOS EL VALOR DE NOMBRE
	pelicula.Director = r.FormValue("director")          //INSERTAMOS EL VALOR DEL CAMPO DIRECTOR
	pelicula.Anio, _ = strconv.Atoi(r.FormValue("anio")) // INSERTAR EL VAMOR DE ANIO CONVERTIDO A INT

	newID, err := p.repo.Create(r.Context(), &pelicula) //SE LLAMA EL METODO DE LA BD Y SE VERIFICA SI SE REALIZO CORRECTAMENTE

	fmt.Println(newID) //IMPRIMIMOS EL NUEVO ID

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "Successfully Created"})
}

// ACTUALIZA LA PELICULA DESDE EL ID EN LA URL
func (p *Pelicula) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id")) //SE OBTIENE EL ID DEL URL Y SE CONVIERTE
	data := models.Pelicula{ID: int(id)}         // ASIGNAMIOS EL ID AL MODELO

	r.ParseForm()                       //VOLVEMOS LA PETICION EN UN dFORMULARIO
	data.Nombre = r.FormValue("nombre") //AGREGAMOS EL VALOR PARA CADA CAMPO
	data.Director = r.FormValue("director")
	data.Anio, _ = strconv.Atoi(r.FormValue("anio"))

	payload, err := p.repo.Update(r.Context(), &data)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// RETORNAR EL DETALLE DE LA PELICULA
func (p *Pelicula) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	payload, err := p.repo.GetByID(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusNoContent, "Content not found")
		return
	}

	respondwithJSON(w, http.StatusOK, payload)
}

// ELIMINAR PELICULA
func (p *Pelicula) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := p.repo.Delete(r.Context(), int64(id))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server Error")
	}

	respondwithJSON(w, http.StatusMovedPermanently, map[string]string{"message": "Delete Successfully"})
}

// METODO QUE PERMITE CONVERTIR LA PETICION EN JSON
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//METODO QUE PERMITE  EL RETORNO DE LOS ERRORES
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondwithJSON(w, code, map[string]string{"message": msg})
}
