package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	"gopkg.in/mgo.v2"
)

var (
	renderer     *render.Render
	mongoSession *mgo.Session
)

func init() {
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	mongoSession = s
}

func main() {
	// create router
	router := httprouter.New()

	router.POST("/users", createUser)
	router.GET("/users", retrieveUsers)

	// negroni 미들웨어 생성
	n := negroni.Classic()
	// negroni에 router를 핸들러로 등록
	n.UseHandler(router)

	// execute web server
	n.Run(":3000")
}
