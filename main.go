package main

import (
	"os"
    "log"
    "net/http"
    "github.com/wkless/ctd/routes"
)

func main() {
    router := routes.NewRouter()

    var listen string = os.Getenv("LISTEN")
    if listen == "" {
    	listen = ":8080"
    }

    log.Println("listening on", listen)
    log.Fatal(http.ListenAndServe(listen, router))
}