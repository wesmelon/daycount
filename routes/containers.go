package routes

import (
	"log"
	"strconv"
    "net/http"
    "encoding/json"
    "github.com/wkless/ctd/models"
    "github.com/gorilla/mux"
)

func getContainer(w http.ResponseWriter, r *http.Request, id string) (models.Container, bool) {
	var c models.Container
	db.QueryRow("SELECT id, uid, cid, name, description, time, is_public FROM containers WHERE id = $1", 
		id).Scan(&c.Id, &c.UserId, &c.CategoryId, &c.Name, &c.Description, &c.Time, &c.IsPublic)

	auth := isUserAuthorized(w, r, c.UserId)
	return c, auth
}

func GetContainerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	
	if c, auth := getContainer(w, r, vars["id"]); auth { //session check
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(c); err != nil {
	       	log.Fatal(err)
	    }
	}
}

func GetContainersByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var c models.Container

	w.Header().Set("Content-Type", "application/json")

	uid, err := strconv.Atoi(vars["uid"])
	if err != nil {
		log.Fatal(err)
	}

	if auth := isUserAuthorized(w, r, uid); auth {
		rows, err := db.Query("SELECT id, uid, cid, name, description, time, is_public FROM containers WHERE uid = $1", vars["uid"])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var containers models.Containers
		for rows.Next() {
			err := rows.Scan(&c.Id, &c.UserId, &c.CategoryId, &c.Name, &c.Description, &c.Time, &c.IsPublic)
			if err != nil {
				log.Fatal(err)
			}

			containers = append(containers, c)
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(containers); err != nil {
	       	log.Fatal(err)
	    }

	    err = rows.Err()

		if err != nil {
	    	log.Fatal(err)
	    }
	}
}

func GetContainersByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var c models.Container

	w.Header().Set("Content-Type", "application/json")

	uid, err := strconv.Atoi(vars["cid"])
	if err != nil {
		log.Fatal(err)
	}

	if auth := isUserAuthorized(w, r, uid); auth {
		rows, err := db.Query("SELECT id, uid, cid, name, description, time, is_public FROM containers WHERE cid = $1", vars["uid"])
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var containers models.Containers
		for rows.Next() {
			err := rows.Scan(&c.Id, &c.UserId, &c.CategoryId, &c.Name, &c.Description, &c.Time, &c.IsPublic)
			if err != nil {
				log.Fatal(err)
			}

			containers = append(containers, c)
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(containers); err != nil {
	       	log.Fatal(err)
	    }

	    err = rows.Err()

		if err != nil {
	    	log.Fatal(err)
	    }
	}
}

func PostContainer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var container models.Container

	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&container); err != nil {
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatal(err)
        }
    }

    if auth := isUserAuthorized(w, r, container.UserId); auth {
		err := db.QueryRow("INSERT INTO containers(uid, cid, name, description, time, is_public) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			container.UserId, container.CategoryId, container.Name, container.Description, container.Time, container.IsPublic).Scan(&container.Id)
		if err != nil {
	    	log.Fatal(err)
	    }

		w.WriteHeader(http.StatusCreated)
	    if err := json.NewEncoder(w).Encode(container); err != nil {
	        log.Fatal(err)
	    }
	}
}

func PutContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	var container models.Container

	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&container); err != nil {
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatal(err)
        }
    }

	if _, auth := getContainer(w, r, vars["id"]); auth { //session check
		err := db.QueryRow("UPDATE containers SET uid = $1, cid = $2, name = $3, description = $4, time = $5, is_public = $6 WHERE id = $7 RETURNING id", 
			container.UserId, container.CategoryId, container.Name, container.Description, container.Time, container.IsPublic, vars["id"]).Scan(&container.Id)
	    if err != nil {
			log.Fatal(err)
	    }

		w.WriteHeader(http.StatusAccepted)
	    if err := json.NewEncoder(w).Encode(container); err != nil {
	        log.Fatal(err)
	    }
	}
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if c, auth := getContainer(w, r, vars["id"]); auth { //session check
	    _, err := db.Exec("DELETE FROM containers WHERE id = $1", vars["id"])
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