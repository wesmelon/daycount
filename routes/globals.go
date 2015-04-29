package routes

import (
	"log"
    "database/sql"
    "github.com/gorilla/mux"
)

const src = "C:/Users/Wesley/Desktop/ctd/src/github.com/wkless/ctd/templates/"
var db* sql.DB

func InitDb() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=test dbname=countingdays sslmode=disable") // need to change SSL
    if err != nil {
        log.Fatal(err)
    }
}

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            HandlerFunc(route.HandlerFunc)
    }

    return router
}