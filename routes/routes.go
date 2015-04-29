package routes

import (
    "path"
    "log"
    "net/http"
    "html/template"
)

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles(path.Join(src, "index.html")))

    err := t.Execute(w, nil)
    if err != nil {
        log.Fatal(err)
    }
} 
 
func InternalPageHandler(w http.ResponseWriter, r *http.Request) {
    id, email := getUserName(r)
    if email != "" {
        t := template.Must(template.ParseFiles(path.Join(src, "internal.html")))

        err := t.Execute(w, map [string] string {"Id": id, "Email": email})
        if err != nil {
            log.Fatal(err)
        }
    } else {
        http.Redirect(w, r, "/", 302)
    }
}

func authHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if auth := isUserAuthorized(w, r, 0); auth {
            fn(w, r)
        }
    }
}
