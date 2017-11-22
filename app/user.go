package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	ID        int8    `json:"id"`
	UserName  string  `json:"user_name"`
	CreatedAt []uint8 `json:"created_at"`
}

func createUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if err := req.ParseForm(); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	stmt, err := db.Prepare("INSERT INTO users (name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(req.FormValue("user_id"))
	if err != nil {
		log.Fatal(err)
	}
}

func retrieveUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rows, err := db.Query("SELECT id, name, created_at FROM users")
	if err != nil {
		panic(err.Error())
	}

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.UserName, &u.CreatedAt); err != nil {
			log.Fatal(err)
			return
		}
		users = append(users, u)
	}
	renderer.JSON(w, http.StatusOK, users)
}
