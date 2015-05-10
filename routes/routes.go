package routes

import (
    "log"
    "net/http"
    "html/template"
    "github.com/gorilla/mux"
)

var tmp *template.Template

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := getUserName(r)
    if id != "" {
        InternalPageHandler(w, r)
    } else {
        t := template.Must(template.ParseFiles(tmpl + "index.html"))

        err := t.Execute(w, nil)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles(tmpl + "auth/signup.html"))

    err := t.Execute(w, nil)
    if err != nil {
        log.Fatal(err)
    }
}
 
func InternalPageHandler(w http.ResponseWriter, r *http.Request) {
    id, email := getUserName(r)
    t := template.Must(template.ParseFiles(tmpl + "internal.html"))

    err := t.Execute(w, map [string] string {"Id": id, "Email": email})
    if err != nil {
        log.Fatal(err)
    }
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
    id, email := getUserName(r)
    tmp = template.New("homepage.html").Delims("<<", ">>")
    t, _ := tmp.ParseFiles(tmpl + "homepage.html")

    err := t.Execute(w, map [string] string {"Id": id, "Email": email})
    if err != nil {
        log.Fatal(err)
    }   
}

func ContainerPageHandler(w http.ResponseWriter, r *http.Request) {
    id, email := getUserName(r)
    vars := mux.Vars(r)

    tmp = template.New("container.html").Delims("<<", ">>")
    t, _ := tmp.ParseFiles(tmpl + "container.html")

    err := t.Execute(w, map [string] string {"Cid": vars["id"], "Id": id, "Email": email})
    if err != nil {
        log.Fatal(err)
    }   
}

