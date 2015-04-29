package routes

import (
	"log"
	"time"
	"strconv"
    "net/http"
    "encoding/json"
    "github.com/wkless/ctd/models"
    "github.com/gorilla/mux"
)

func getCategory(w http.ResponseWriter, r *http.Request, id string) (models.Category, bool) {
	var c models.Category
	db.QueryRow("SELECT id, uid, name, picture, creation_time FROM categories WHERE id = $1", 
		id).Scan(&c.Id, &c.UserId, &c.Name, &c.Picture, &c.CreatedTime)

	auth := isUserAuthorized(w, r, c.UserId)
	return c, auth
}

func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	
	if c, auth := getCategory(w, r, vars["id"]); auth { //session check
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(c); err != nil {
	       	log.Fatal(err)
	    }
	}
}

func GetCategoriesByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var c models.Category

	w.Header().Set("Content-Type", "application/json")

	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		log.Fatal(err)
	}

	if auth := isUserAuthorized(w, r, uid); auth {
		rows, err := db.Query("SELECT id, uid, name, picture, creation_time FROM categories WHERE uid = $1", vars["uid"])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var categories models.Categories
		for rows.Next() {
			err := rows.Scan(&c.Id, &c.UserId, &c.Name, &c.Picture, &c.CreatedTime)
			if err != nil {
				log.Fatal(err)
			}

			categories = append(categories, c)
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(categories); err != nil {
	       	log.Fatal(err)
	    }

	    err = rows.Err()

		if err != nil {
	    	log.Fatal(err)
	    }
	}
}

func PostCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var category models.Category

	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&category); err != nil {
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatal(err)
        }
    }

    if auth := isUserAuthorized(w, r, category.UserId); auth {
	    category.CreatedTime = time.Now()
		err := db.QueryRow("INSERT INTO categories(uid, name, picture, creation_time) VALUES ($1, $2, $3, $4) RETURNING id",
			category.UserId, category.Name, category.Picture, category.CreatedTime).Scan(&category.Id)
		if err != nil {
	    	log.Fatal(err)
	    }

		w.WriteHeader(http.StatusCreated)
	    if err := json.NewEncoder(w).Encode(category); err != nil {
	        log.Fatal(err)
	    }
	}
}

func PutCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	
	var category models.Category
	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&category); err != nil {
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatal(err)
        }
    }

	if _, auth := getCategory(w, r, vars["id"]); auth { //session check
		err := db.QueryRow("UPDATE categories SET uid = $1, name = $2, picture = $3, creation_time = $4 WHERE id = $5 RETURNING id", 
			category.UserId, category.Name, category.Picture, category.CreatedTime, vars["id"]).Scan(&category.Id)
	    if err != nil {
			log.Fatal(err)
	    }

		w.WriteHeader(http.StatusAccepted)
	    if err := json.NewEncoder(w).Encode(category); err != nil {
	        log.Fatal(err)
	    }
	}
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if c, auth := getCategory(w, r, vars["id"]); auth { //session check
	    _, err := db.Exec("DELETE FROM categories WHERE id = $1", vars["id"])
	    if err != nil {
			log.Fatal(err)
	    }
	    w.WriteHeader(http.StatusAccepted)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(c); err != nil {
	       	log.Fatal(err)
	    }
	}
}