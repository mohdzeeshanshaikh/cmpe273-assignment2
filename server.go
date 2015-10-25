package main

import (
    "fmt"
    "log"
    "gopkg.in/mgo.v2"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "./controller"
)

func getMongoSession() *mgo.Session {
    session, err := mgo.Dial("mongodb://mohdzeeshanshaikh:password@ds045454.mongolab.com:45454/cmpe273assignment2")
    if err != nil {
        panic(err)
    }
    return session
}

func main() {
    router := httprouter.New()
    friendController := controller.NewFriendController(getMongoSession())
    router.POST("/locations", friendController.CreateFriend)
    router.GET("/locations/:id", friendController.GetFriend)
    router.PUT("/locations/:id", friendController.UpdateFriend)
    router.DELETE("/locations/:id", friendController.RemoveFriend)
    fmt.Println("Server listening on port 8080")
	  log.Fatal(http.ListenAndServe(":8080", router))
}
