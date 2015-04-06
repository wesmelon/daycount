package main

import (
    "log"
    "fmt"
    "net/http"
    "github.com/wkless/ctd/routes"
)

func main() {
    router := routes.NewRouter()

    fmt.Printf("listening 8080\n")
    log.Fatal(http.ListenAndServe(":8080", router))
}