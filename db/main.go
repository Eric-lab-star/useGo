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
	//open pgx driver
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Get("/", home(nil))
	r.Post("/", home(db))
	r.Get("/album", album(db))
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	fileserver(r, "/assets", filesDir) // serve filesDir files from /assets
	http.ListenAndServe("localhost:3000", r)

}

type userData struct {
	Uuid      uuid.UUID
	Firstname string
	Lastname  string
	Date      string
	Doc       string
	Imgurl    string
}

func home(db *sql.DB) http.HandlerFunc {

	templ, error := template.ParseFiles("index.html")
	if error != nil {
		panic("failed to parse template")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := templ.Execute(w, nil)
			if err != nil {
				panic(err)
			}
		}

		if r.Method == "POST" {
			userdata, err := FileSave(w, r)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = insert(db, userdata)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			w.Write([]byte("saved file"))
		}
	}
}

func album(db *sql.DB) http.HandlerFunc {
	rows, err := db.Query(`
		SELECT * FROM users ;
	`,
	)
	if err != nil {
		log.Printf("Failed to query")
	}
	users := []userData{}
	user := userData{}
	for rows.Next() {
		err := rows.Scan(&user.Uuid, &user.Firstname, &user.Lastname, &user.Date, &user.Doc, &user.Imgurl)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("album.html")
		if err != nil {
			panic("failed to parse album.html")
		}
		err = temp.Execute(w, users)
		if err != nil {
			log.Printf("unable to execute template 'album.html'\n")
			os.Exit(1)
		}
	}
}

func insert(db *sql.DB, userData *userData) error {
	_, err := db.Exec(
		`
	INSERT INTO users (id, firstname, lastname,date, documents, imgurl) 
	VALUES(gen_random_uuid (), $1, $2, current_date, $3,$4)
	`,
		userData.Firstname, userData.Lastname, userData.Doc, userData.Imgurl)
	if err != nil {
		return err

	}
	return nil
}

func FileSave(w http.ResponseWriter, r *http.Request) (*userData, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatalf("parseMulipartForm err %v", err)
	}
	// Retrieve the file from form data
	imgfile, header, err := r.FormFile("file")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	defer imgfile.Close()
	// datas to be
	title := r.Form.Get("title") + filepath.Ext(header.Filename)
	firstname := r.Form.Get("firstname")
	lastname := r.Form.Get("lastname")
	doc := r.Form.Get("doc")
	imgUrl := filepath.Join("assets", "uploaded", title)

	wd, err := os.Getwd()
	if err != nil {
		log.Println("unable to get working directory")
		return nil, err
	}
	dirpath := filepath.Join(wd, "assets", "uploaded")
	err = os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Println("unable to mkdir")
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(dirpath, title), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Printf("openFile err: %v \n", err)
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, imgfile)
	if err != nil {
		log.Printf("Copy file err: %v \n", err)
		return nil, err
	}

	return &userData{
		Firstname: firstname,
		Lastname:  lastname,
		Doc:       doc,
		Imgurl:    imgUrl,
	}, nil
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
