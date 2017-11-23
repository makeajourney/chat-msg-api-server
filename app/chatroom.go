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
	renderer.JSON(w, http.StatusCreated, findRoom(string(chatroomID)))
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

// /chatrooms/:room_id/user/:user_id/:status
func joinChatroom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	status := ps.ByName("status")
	roomID := ps.ByName("room_id")
	userID := ps.ByName("user_id")

	var userStatus int
	switch status {
	case "join":
		userStatus = 1
	case "leave":
		userStatus = 2
	default:
		http.Error(w, "status '"+status+"' is not supported", http.StatusNotFound)
	}
	stmt, err := db.Prepare("UPDATE chatrooms_users SET status=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(userStatus, roomID, userID)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, findRoom(roomID))
}

// "/chatrooms/:room_id/open", openChatroom
func openChatroom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	roomID := ps.ByName("room_id")

	userStatus := 1
	stmt, err := db.Prepare("UPDATE chatrooms_users SET user_status=? WHERE chatroom_id=?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(userStatus, roomID)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, findRoom(roomID))
}

// query
func findRoom(roomID string) *Chatroom {
	rows, err := db.Query(
		"select A.chatroom_id, A.user_id, A.user_status, B.user_id, B.user_status from chatrooms_users as A inner join chatrooms_users as B where A.chatroom_id = B.chatroom_id and A.chatroom_id=" + roomID,
	)
	if err != nil {
		log.Fatal(err)
	}

	cr := new(Chatroom)
	for rows.Next() {
		var u1 ChatUser
		var u2 ChatUser
		if err := rows.Scan(&cr.ID, &u1.UserID, &u1.Status, &u2.UserID, &u2.Status); err != nil {
			log.Fatal(err)
		}
		cr.CurrentUser = u1
		cr.Partner = u2
	}
	return cr
}
