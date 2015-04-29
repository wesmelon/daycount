package routes

import (
    "log"
    "net/http"
    "strconv"
)

func getUserName(request *http.Request) (id string, email string) {
    if cookie, err := request.Cookie("session"); err == nil {
        value := make(map[string]string)
        if err = CookieHandler.Decode("session", cookie.Value, &value); err == nil {
            id = value["id"]
            email = value["email"]
        }
    }
    return id, email
}

func isUserAuthorized(w http.ResponseWriter, r *http.Request, uid int) bool {
    if cookie, err := r.Cookie("session"); err == nil {
        value := make(map[string]string)
        if err = CookieHandler.Decode("session", cookie.Value, &value); err == nil {
            if uid != 0 {
                sessionId, err := strconv.Atoi(value["id"])
                if err != nil {
                    log.Fatal(err)
                }
                if sessionId != uid {
                    http.Redirect(w, r, "/", http.StatusUnauthorized)
                    return false
                }
            }
            return true
        }
    }

    http.Redirect(w, r, "/", http.StatusUnauthorized)
    return false
}
