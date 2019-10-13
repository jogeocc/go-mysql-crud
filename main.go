package main

import (
	"fmt"
	"net/http"
	"os"

	ph "./handler/http"

	"./driver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load() //CARGAR ARCHIVO .env con las variables que se usaran dentro del sistema
}

func main() {

	// ---------------------- CARGANDO LOS VALORES DE LA BD DESDE EL ARCHIVO ENV -------------------------
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	//-----------------------------------------------------------------------------------------------------

	//_______________________________ ABRIMOS CONEXION DE LA BD ___________________________________________
	connection, err := driver.ConnectSQL(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	//______________________________________________________________________________________________________

	//---------------------------- DANDO DE ALTA LA RUTA PRINCIPAL Y DAR DE ALTA EL MIDDLEWARE ------------
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	//------------------------------------------------------------------------------------------------------

	//_________________DAR DE ALTA LOS METODOS Y LAS RUTAS DEL MODULO_____________________________________
	pHandler := ph.NewPeliculaHandler(connection)
	r.Route("/", func(rt chi.Router) { // DEFINIMOS RUTA BASE
		rt.Mount("/peliculas", peliculaRouter(pHandler)) //MONTAMOS LA RUTA DE PELICULAS Y SUS METODOS
	})

	//_____________________________________________________________________________________________________

	//-------------------- ARRANCAR EL SERVIDOR EN EL PUERTO DESEADO ---------------------------------------
	fmt.Println("Server listen at :8005")
	http.ListenAndServe(":8005", r) // PUERTO Y SERVIDOR QUE SE INICIAN
	//------------------------------------------------------------------------------------------------------
}

// Definicion de las rutas y el middleware para validar los datos de la URL
func peliculaRouter(pHandler *ph.Pelicula) http.Handler {
	r := chi.NewRouter()
	r.Get("/", pHandler.Fetch)
	r.Get("/{id:[0-9]+}", pHandler.GetByID)
	r.Post("/", pHandler.Create)
	r.Put("/{id:[0-9]+}", pHandler.Update)
	r.Delete("/{id:[0-9]+}", pHandler.Delete)
	return r
}
