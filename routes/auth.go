package routes

import (
    "fmt"
    "bytes"
    "time"
    "log"
    "net/http"
    "strconv"
    "encoding/base64"
    "crypto/rand"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
    "github.com/wkless/ctd/models"
    "github.com/gorilla/securecookie"
)

const SaltSize = 16

var CookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
    securecookie.GenerateRandomKey(32))

func LoginHandler(response http.ResponseWriter, request *http.Request) {
     email := request.FormValue("email")
     pass := request.FormValue("password")
     redirectTarget := "/"
     if email != "" && pass != "" {
        authenticated := authenticatePassword(email, pass, response)
        if authenticated {
            redirectTarget = "/internal"
        } else {
            redirectTarget = "/"
        }
    }
    http.Redirect(response, request, redirectTarget, 302)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
    clearSession(response)
    http.Redirect(response, request, "/", 302)
}

func setSession(id, email string, response http.ResponseWriter) {
    value := map[string]string{
        "id": string(id),
        "email": email,
    }
    if encoded, err := CookieHandler.Encode("session", value); err == nil {
        cookie := &http.Cookie{
            Name:  "session",
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(response, cookie)
    }
}

func clearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "session",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}

func authenticatePassword(email, password string, response http.ResponseWriter) bool {
    var id int
    var dbUserEmail, dbPasswordHash, dbPasswordSalt string
    err := db.QueryRow("SELECT id, email, password_hash, password_salt FROM users WHERE email = $1", 
        email).Scan(&id, &dbUserEmail, &dbPasswordHash, &dbPasswordSalt)
    if err != nil {
        return false
    }

    err = bcrypt.CompareHashAndPassword([]byte(dbPasswordHash), []byte(dbPasswordSalt+password))
    if err != nil {
        return false
    }

    setSession(strconv.Itoa(id), dbUserEmail, response)
    return true
}

const signupPage = `
<h1>Signup</h1>
<form method="post" action="/signup">
    <label for="Email">Email</label>
    <input type="Email" id="Email" name="Email">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <label for="Fullname">Fullname</label>
    <input type="Fullname" id="Fullname" name="Fullname">
    <button type="submit">Login</button>
</form>`

func SignupHandler(response http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(response, signupPage)
}

func GoSignup(response http.ResponseWriter, request *http.Request) {
     email := request.FormValue("Email")
     pass := request.FormValue("password")
     fullname := request.FormValue("Fullname")
     redirectTarget := "/"
     if email != "" && pass != "" && fullname != "" {
        hash, salt := hashPassword(pass);

        user := models.User {
            "needid", 
            email,
            fullname,
            hash,
            salt,
            time.Now(),
            time.Now(),
            false,
            false,
        }

        storeUser(user)

        redirectTarget = "/"
    }
    http.Redirect(response, request, redirectTarget, 302)
}

func hashPassword(password string) (string, string) {
    var buffer bytes.Buffer

    b := make([]byte, SaltSize)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }

    salt := base64.URLEncoding.EncodeToString(b)

    buffer.WriteString(salt)
    buffer.WriteString(password)

    hash, err := bcrypt.GenerateFromPassword(buffer.Bytes(), 10)
    if err != nil {
        log.Fatal(err)
    }

    return string(hash), salt
}

func storeUser(user models.User) {
    const layout = "2006-01-02 15:04:05"

    _, err := db.Exec("INSERT INTO users(email, full_name, password_hash, password_salt, is_disabled, is_activated, creation_time, last_login_time) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", 
                            user.Email, user.FullName, user.PasswordHash, user.PasswordSalt, strconv.FormatBool(user.IsDisabled),
                            strconv.FormatBool(user.IsActivated), user.CreatedTime.Format(layout), user.LastLoginTime.Format(layout))
    if err != nil {
        log.Fatal(err)
    }
}