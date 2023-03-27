package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var firstname string
	var uuid uuid.UUID
	var lastname string
	err = db.QueryRow("select id, firstname, lastname from users").Scan(&uuid, &firstname, &lastname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows, err := db.Query("select * from users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)

	r := chi.NewRouter()
	r.Get("/", home)
	r.Post("/", home)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	fileserver(r, "/files", filesDir)
	http.ListenAndServe("localhost:3000", r)

}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		templ, error := template.ParseFiles("index.html")
		if error != nil {
			panic("failed to parse template")
		}
		err := templ.Execute(w, "http://localhost:3000/files/moon.jpg")
		if err != nil {
			panic(err)
		}
	}
	if r.Method == "POST" {
		path := FileSave(r)
		fmt.Fprintf(w, "uploaded result:%v \n", path)

		fmt.Println(path)

	}
}

func FileSave(r *http.Request) string {
	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatalf("parsMulipartForm err %v", err)
	}
	title := r.Form.Get("title")
	desc := r.Form.Get("desc")

	fmt.Println(title)
	fmt.Println(desc)
	// Retrieve the file from form data
	f, h, err := r.FormFile("file")
	if err != nil {
		log.Fatalf("FormFile err: %v \n", err)
	}
	defer f.Close()
	wd, err := os.Getwd()
	path := filepath.Join(wd, "uploaded")
	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/" + title + filepath.Ext(h.Filename)
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("openFile err: %v \n", err)
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		log.Fatalf("Copy file err: %v \n", err)
	}
	return title + filepath.Ext(h.Filename)
}

func fileserver(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
