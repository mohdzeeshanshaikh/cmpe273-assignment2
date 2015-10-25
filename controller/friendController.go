package controller

import (
    "../model"
    "fmt"
    "net/http"
    "net/url"
    "encoding/json"
    "strings"
    "errors"
    "log"
    "io"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/julienschmidt/httprouter"
)

type FriendController struct{
    session *mgo.Session
}

func NewFriendController(s *mgo.Session) *FriendController {
    return &FriendController{s}
}

func fetchCoordinates(Friend *model.Friend) error {
client := &http.Client{}
	address := Friend.Address + "+" + Friend.City + "+" + Friend.State + "+" + Friend.Zip;

    mapsUrl := "http://maps.google.com/maps/api/geocode/json?address="

	mapsUrl += url.QueryEscape(address)
	mapsUrl += "&sensor=false"

	req, err := http.NewRequest("GET", mapsUrl , nil)
	res, err := client.Do(req)

    if err != nil {
        return err
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return err
    }
    defer res.Body.Close()

    var contents map[string]interface{}
    err = json.Unmarshal(body, &contents)
    if err != nil {
        return err
    }

    if !strings.EqualFold(contents["status"].(string), "OK") {
        return errors.New("Coordinates unavailable")
    }

    results := contents["results"].([]interface{})
    location := results[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"]

    Friend.Coordinate.Lat = location.(map[string]interface{})["lat"].(float64)
    Friend.Coordinate.Lng = location.(map[string]interface{})["lng"].(float64)

    if err != nil {
        return err
    }

    return nil
}

func (cc FriendController) CreateFriend(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
    Friend := model.Friend{}
    json.NewDecoder(req.Body).Decode(&Friend)
    Friend.Id = bson.NewObjectId()
    err := fetchCoordinates(&Friend)
    if err != nil {
		log.Println(err)
	}

    conn := cc.session.DB("cmpe273assignment2").C("users")
    err = conn.Insert(Friend)
    if err != nil {
        log.Println(err)
    }

    friendJson, _ := json.Marshal(Friend)
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", friendJson)
    }
}

func (cc FriendController) GetFriend(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    friend, err := fetchFriendById(cc, param.ByName("id"))
    if err != nil {
		log.Println(err)
	}

    friendJson, _ := json.Marshal(friend)
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", friendJson)
    }
}

func (cc FriendController) UpdateFriend(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
    updatedUsr, err := updateFriendLocation(cc, param.ByName("id"), req.Body)
    if err != nil {
		log.Println(err)
	}

    friendJson, _ := json.Marshal(updatedUsr)
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", friendJson)
    }
}

func (cc FriendController) RemoveFriend(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {

    friend, err := fetchFriendById(cc, param.ByName("id"))
    if err != nil {
		log.Println(err)
        log.Println(friend)
	}

    objId := bson.ObjectIdHex(param.ByName("id"))
    conn := cc.session.DB("cmpe273assignment2").C("users")
    err = conn.Remove(bson.M{"id": objId})
    if err != nil {
        log.Println(err)
    }
    rw.Header().Set("Content-Type", "plain/text")
    if err != nil {
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "Friend ID=%s has been deleted", param.ByName("id"))
    }
}

func fetchFriendById(cc FriendController, id string) (model.Friend, error) {

    if !bson.IsObjectIdHex(id) {
        return model.Friend{}, errors.New("Invalid Friend ID")
    }
    objId := bson.ObjectIdHex(id)
    Friend := model.Friend{}
    conn := cc.session.DB("cmpe273assignment2").C("users")
    err := conn.Find(bson.M{"id": objId}).One(&Friend)
    if err != nil {
        return model.Friend{}, errors.New("This Friend Id doesn't exists")
    }
    return Friend, nil
}

func updateFriendLocation(cc FriendController, id string, contents io.Reader) (model.Friend, error) {

    friend, err := fetchFriendById(cc, id)
    if err != nil {
        return model.Friend{}, err
    }

    updFriend := model.Friend{}
    updFriend.Id = friend.Id
    updFriend.Name = friend.Name
    json.NewDecoder(contents).Decode(&updFriend)

    err = fetchCoordinates(&updFriend)
    if err != nil {
        return model.Friend{}, err
    }

    objId := bson.ObjectIdHex(id)
    conn := cc.session.DB("cmpe273assignment2").C("users")
    err = conn.Update(bson.M{"id": objId}, updFriend)
    if err != nil {
        log.Println(err)
        return model.Friend{}, errors.New("Given id is invalid")
    }
    return updFriend, nil
}
