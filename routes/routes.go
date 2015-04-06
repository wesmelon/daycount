package routes

import (
    "fmt"
    "net/http"
    "log"
    "strconv"
    "database/sql"
    _ "github.com/lib/pq"

    "github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
    db, err = sql.Open("postgres", "user=postgres password=test dbname=countingdays sslmode=disable") // need to change SSL
    if err != nil {
        log.Fatal(err)
    }

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

const indexPage = `
<h1>Count The Days...</h1>
<form method="post" action="/login">
    <label for="email">Email</label>
    <input type="text" id="email" name="email">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
    <p><a href="/signup">Signup</a></p>
</form>`

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s, id: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(response, indexPage)
} 
 
func InternalPageHandler(response http.ResponseWriter, request *http.Request) {
    id, email := getUserName(request)
    if email != "" {
        fmt.Fprintf(response, internalPage, email, id)
    } else {
        http.Redirect(response, request, "/", 302)
    }
}

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
                    panic(err)
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

func authHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if auth := isUserAuthorized(w, r, 0); auth {
            fn(w, r)
        }
    }
}

var routes = Routes{
    Route{"Index", "GET", "/", IndexPageHandler},
    Route{"Internal", "GET", "/internal", InternalPageHandler}, 

    Route{"Login", "POST", "/login", LoginHandler}, 
    Route{"Logout", "POST", "/logout", LogoutHandler}, 
    Route{"Signup", "GET", "/signup", SignupHandler}, 
    Route{"GoSignup", "POST", "/signup", GoSignup}, 

    Route{"TodoIndex", "GET", "/todos", TodoIndex}, 
    Route{"TodoShow", "GET", "/todos/{todoId}", TodoShow}, 

    Route{"GetCategoriesByUser", "GET", "/api/categories/user/{uid}", authHandler(GetCategoriesByUser)},
    Route{"GetCategoryById", "GET", "/api/categories/{id}", authHandler(GetCategoryById)}, 
    Route{"PostCategory", "POST", "/api/categories/", authHandler(PostCategory)}, 
    Route{"PutCategory", "PUT", "/api/categories/{id}", authHandler(PutCategory)}, 
    Route{"DeleteCategory", "DELETE", "/api/categories/{id}", authHandler(DeleteCategory)}, 

    Route{"GetContainersByUser", "GET", "/api/containers/user/{uid}", authHandler(GetContainersByUser)},
    Route{"GetContainerById", "GET", "/api/containers/{id}", authHandler(GetContainerById)}, 
    Route{"PostContainer", "POST", "/api/containers/", authHandler(PostContainer)},
    Route{"PutContainer", "PUT", "/api/containers/{id}", authHandler(PutContainer)},  
    Route{"DeleteContainer", "DELETE", "/api/containers/{id}", authHandler(DeleteContainer)}, 

    Route{"GetDatesByContainer", "GET", "/api/dates/container/{cid}", authHandler(GetDatesByContainer)},
    Route{"GetDateById", "GET", "/api/dates/{id}", authHandler(GetDateById)}, 
    Route{"PostDate", "POST", "/api/dates/", authHandler(PostDate)},
    Route{"PutDate", "PUT", "/api/dates/{id}", authHandler(PutDate)},  
    Route{"DeleteDate", "DELETE", "/api/dates/{id}", authHandler(DeleteDate)}, 
}