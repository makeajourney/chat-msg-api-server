package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	UserId string        `bson:"user_id" json:"user_id"`
}

func (u *User) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&u.UserId: "user_id"}
}

func createUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	u := new(User)
	errs := binding.Bind(req, u)
	// if errs.Handle(w) {
	// 	return
	// }
	if errs != nil {
		return
	}

	session := mongoSession.Copy()
	defer session.Close()

	u.ID = bson.NewObjectId()
	c := session.DB("poc").C("users")

	if err := c.Insert(u); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusCreated, u)
}

func retrieveUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	var users []User
	err := session.DB("poc").C("users").Find(nil).All(&users)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, users)
}
