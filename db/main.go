package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctxbg := context.Background()
	conn, err := pgx.Connect(ctxbg, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctxbg)

	var firstname string
	var uuid uuid.UUID
	var lastname string
	err = conn.QueryRow(ctxbg, "select id, firstname, lastname from users").Scan(&uuid, &firstname, &lastname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", homeHandler(firstname, lastname, uuid))
	http.ListenAndServe("localhost:3000", nil)

	fmt.Println(uuid, firstname, lastname)
}

func homeHandler(firstname string, lastname string, id uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s %s, %v", firstname, lastname, id)
	}
}
