package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	renderer *render.Render
	db       *sql.DB
	err      error
)

func init() {
	renderer = render.New()
}

func main() {
	db, err = sql.Open("mysql", "root:root@/chatpoc")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// create router
	router := httprouter.New()

	router.POST("/users", createUser)
	router.GET("/users", retrieveUsers)
	router.GET("/users/:user_id/chatrooms", findRoomByUser)

	router.POST("/chatrooms", createRoom)
	router.GET("/chatrooms/:room_id/open", openChatroom)
	router.GET("/chatrooms/:room_id/users/:user_id/:status", joinChatroom)

	// negroni 미들웨어 생성
	n := negroni.Classic()
	// negroni에 router를 핸들러로 등록
	n.UseHandler(router)

	// execute web server
	n.Run(":3000")
}
