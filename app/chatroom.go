package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Chatroom struct {
	ID          int8     `json:"id"`
	CurrentUser ChatUser `json:"user"`
	Partner     ChatUser `json:"partner"`
	CreatedAt   []uint8  `json:"created_at"`
}

type ChatUser struct {
	UserID int8 `json:"user_id"`
	Status int  `json:"status"`
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if err := req.ParseForm(); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	user1 := req.PostForm.Get("user1")
	user2 := req.PostForm.Get("user2")

	stmt, err := db.Prepare("INSERT INTO chatrooms () VALUES ()")
	if err != nil {
		log.Fatal(err)
	}
	rst, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	chatroomID, err := rst.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err = db.Prepare("INSERT INTO chatrooms_users (chatroom_id, user_id, user_status) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(chatroomID, user1, false)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(chatroomID, user2, false)
	if err != nil {
		log.Fatal(err)
	}
	renderer.JSON(w, http.StatusCreated, chatroomID)
}

func findRoomByUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user_id")

	rows, err := db.Query(
		"select A.chatroom_id, A.user_id, A.user_status, B.user_id, B.user_status from chatrooms_users as A inner join chatrooms_users as B where A.chatroom_id = B.chatroom_id and A.user_id != B.user_id and A.user_id=" + userID,
	)
	if err != nil {
		log.Fatal(err)
	}

	chatrooms := []Chatroom{}
	for rows.Next() {
		var cr Chatroom
		var u1 ChatUser
		var u2 ChatUser
		if err := rows.Scan(&cr.ID, &u1.UserID, &u1.Status, &u2.UserID, &u2.Status); err != nil {
			log.Fatal(err)
			return
		}
		cr.CurrentUser = u1
		cr.Partner = u2
		chatrooms = append(chatrooms, cr)
	}
	renderer.JSON(w, http.StatusOK, chatrooms)
}
