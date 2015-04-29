package routes

import (
	"log"
    "net/http"
    "encoding/json"
    "github.com/wkless/ctd/models"
    "github.com/gorilla/mux"
)

func getUidOfContainer(cid int) int {
	var uid int
	if err := db.QueryRow("SELECT uid FROM containers WHERE id = $1", cid).Scan(&uid); err != nil {
		log.Fatal(err)
	}
	return uid
}

func getUidOfDate(did string) int {
	var uid int
	if err := db.QueryRow("SELECT uid FROM containers WHERE id = (SELECT cid FROM dates WHERE id = $1)", did).Scan(&uid); err != nil {
		log.Fatal(err)
	}
	return uid
}

func GetDateById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var d models.Date
	w.Header().Set("Content-Type", "application/json")

	if isUserAuthorized(w, r, getUidOfDate(vars["id"])) {
		db.QueryRow("SELECT id, cid, name, type, time, icon, content FROM dates WHERE id = $1", 
			vars["id"]).Scan(&d.Id, &d.ContainerId, &d.Name, &d.Type, &d.Time, &d.Icon, &d.Content)
		
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(d); err != nil {
	       	log.Fatal(err)
	    }
	}
}

func GetDatesByContainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var d models.Date

	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, cid, name, type, time, icon, content FROM dates WHERE cid = $1", vars["cid"])
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var dates models.Dates
	for rows.Next() {
		err := rows.Scan(&d.Id, &d.ContainerId, &d.Name, &d.Type, &d.Time, &d.Icon, &d.Content)
		if err != nil {
			log.Fatal(err)
		}

		dates = append(dates, d)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dates); err != nil {
       	log.Fatal(err)
    }

    err = rows.Err()

	if err != nil {
    	log.Fatal(err)
    }
}

func PostDate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var d models.Date

	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&d); err != nil {
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatal(err)
        }
    }

    if isUserAuthorized(w, r, getUidOfContainer(d.ContainerId)) {
	    err := db.QueryRow("INSERT INTO dates(cid, name, type, time, icon, content) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			d.ContainerId, d.Name, d.Type, d.Time, d.Icon, d.Content).Scan(&d.Id)
		if err != nil {
	    	log.Fatal(err)
	    }

		w.WriteHeader(http.StatusCreated)
	    if err := json.NewEncoder(w).Encode(d); err != nil {
	        log.Fatal(err)
	    }
	}
}


func PutDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	var date models.Date

	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&date); err != nil {
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            log.Fatal(err)
        }
    }

	if isUserAuthorized(w, r, getUidOfDate(vars["id"])) {
		err := db.QueryRow("UPDATE dates SET cid = $1, name = $2, type = $3, time = $4, icon = $5, content = $6 WHERE id = $7 RETURNING id", 
			date.ContainerId, date.Name, date.Type, date.Time, date.Icon, date.Content, vars["id"]).Scan(&date.Id)
	    if err != nil {
			log.Fatal(err)
	    }

		w.WriteHeader(http.StatusAccepted)
	    if err := json.NewEncoder(w).Encode(date); err != nil {
	        log.Fatal(err)
	    }
	}
}


func DeleteDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if isUserAuthorized(w, r, getUidOfDate(vars["id"])) {
	    _, err := db.Exec("DELETE FROM dates WHERE id = $1", vars["id"])
	    if err != nil {
			log.Fatal(err)
	    }
	    w.WriteHeader(http.StatusAccepted)
	}
}