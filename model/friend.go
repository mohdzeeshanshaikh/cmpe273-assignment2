package model

import "gopkg.in/mgo.v2/bson"

type (
    Friend struct {
        Id                           bson.ObjectId  			`json:"id" bson:"id"`
        Name                         string          			`json:"name" bson:"name"`
        Address                      string         			`json:"address" bson:"address"`
        City                         string         			`json:"city" bson:"city"`
        State                        string          			`json:"state" bson:"state"`
        Zip                          string          			`json:"zip" bson:"zip"`
        Coordinate                   coordinate           `json:"coordinate" bson:"coordinate"`
    }
    coordinate struct {
        Lat  						float64 					`json:"lat" bson:"lat"`
        Lng							float64 					`json:"lng" bson:"lng"`
    }
)
