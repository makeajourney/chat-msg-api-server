package main

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

type Chatroom struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Users     []string      `bson:"users" json:"users"`
	Timestamp time.Time     `bson:"ts" json:"ts"`
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if err := req.ParseForm(); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	cr := new(Chatroom)
	users := []string{req.PostForm.Get("user1"), req.PostForm.Get("user2")}

	session := mongoSession.Copy()
	defer session.Close()

	cr.ID = bson.NewObjectId()
	cr.Users = users
	cr.Timestamp = time.Now()
	c := session.DB("poc").C("chatrooms")

	if err := c.Insert(cr); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusCreated, cr)
}

func findRoomByUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user_id")

	session := mongoSession.Copy()
	defer session.Close()

	var chatrooms []Chatroom
	err := session.DB("poc").C("chatrooms").Find(bson.M{"users": userID}).Sort("-ts").All(&chatrooms)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, chatrooms)
}
