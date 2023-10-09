package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "data"
)

var db *sql.DB

type Informacion struct {
	ID      int
	Usuario string
	Pass    string
}

func main() {
	// Configurar la conexión a la base de datos PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear una tabla "informacion" si no existe
	crearTabla()

	// Configurar rutas y manejadores web
	r := mux.NewRouter()
	r.HandleFunc("/", mostrarInformacionWeb).Methods("GET")
	r.HandleFunc("/insertar", insertarInformacionWeb).Methods("POST")
	r.HandleFunc("/actualizar/{id}", actualizarInformacionWeb).Methods("POST")
	r.HandleFunc("/eliminar/{id}", eliminarInformacionWeb).Methods("POST")

	// Servir archivos estáticos (CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configurar servidor web
	http.Handle("/", r)
	fmt.Println("Servidor web en ejecución en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func crearTabla() {
	createTableSQL := `
       CREATE TABLE IF NOT EXISTS informacion (
           id SERIAL PRIMARY KEY,
           usuario VARCHAR(255),
           pass VARCHAR(255)
       );
       `

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func consultarInformacion() []Informacion {
	rows, err := db.Query("SELECT * FROM informacion;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var infos []Informacion
	for rows.Next() {
		var info Informacion
		err := rows.Scan(&info.ID, &info.Usuario, &info.Pass)
		if err != nil {
			log.Fatal(err)
		}
		infos = append(infos, info)
	}
	return infos
}

func mostrarInformacionWeb(w http.ResponseWriter, r *http.Request) {
	infos := consultarInformacion()

	tmpl, err := template.New("index").Parse(`
       <!DOCTYPE html>
       <html>
       <head>
           <title>Mi Aplicación CRUD</title>
       </head>
       <body>
           <h1>CRUD de Información</h1>
           <h2>Lista de Información</h2>
           <ul>
               {{range .}}
                   <li>
                       ID: {{.ID}}, Usuario: {{.Usuario}}, Contraseña: {{.Pass}}
                       <form method="post" action="/actualizar/{{.ID}}">
                           <input type="text" name="nuevaPass" placeholder="Nueva Contraseña">
                           <input type="submit" value="Actualizar">
                       </form>
                       <form method="post" action="/eliminar/{{.ID}}">
                           <input type="submit" value="Eliminar">
                       </form>
                   </li>
               {{end}}
           </ul>
           <h2>Insertar Nueva Información</h2>
           <form method="post" action="/insertar">
               <input type="text" name="usuario" placeholder="Usuario" required>
               <input type="text" name="pass" placeholder="Contraseña" required>
               <input type="submit" value="Insertar">
           </form>
       </body>
       </html>
   `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, infos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func insertarInformacion(usuario, pass string) {
	insertSQL := "INSERT INTO informacion (usuario, pass) VALUES ($1, $2);"
	_, err := db.Exec(insertSQL, usuario, pass)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Información insertada con éxito.")
}

func insertarInformacionWeb(w http.ResponseWriter, r *http.Request) {
	usuario := r.FormValue("usuario")
	pass := r.FormValue("pass")

	insertarInformacion(usuario, pass)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func actualizarInformacion(id int, nuevaPass string) {
	updateSQL := "UPDATE informacion SET pass = $2 WHERE id = $1;"
	_, err := db.Exec(updateSQL, id, nuevaPass)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Contraseña actualizada con éxito.")
}

func actualizarInformacionWeb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	nuevaPass := r.FormValue("nuevaPass")

	// Convierte el id de string a int
	idInt := 0
	fmt.Sscanf(id, "%d", &idInt)

	actualizarInformacion(idInt, nuevaPass)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func eliminarInformacion(id int) {
	deleteSQL := "DELETE FROM informacion WHERE id = $1;"
	_, err := db.Exec(deleteSQL, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Información eliminada con éxito.")
}

func eliminarInformacionWeb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convierte el id de string a int
	idInt := 0
	fmt.Sscanf(id, "%d", &idInt)

	eliminarInformacion(idInt)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
