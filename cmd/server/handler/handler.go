package handler

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	db *mongo.Database
}

type Document struct {
	UserId string `json:"userid" bson:"userid"`
	Name   string `json:"name" bson:"name"`
	Age    int    `json:"age" bson:"age"`
}

type NewCollect struct {
	Collection_name string `bson:"collection_name"`
}

type Doc struct {
	Document Document `bson:"document"`
}

func NewDocumentHandler(c *mongo.Client) *Handler {
	db := c.Database("app")
	return &Handler{db: db}
}

type PutReqBodyDoc struct {
	UserId string `json:"userid" bson:"userid"`
	Name   string `json:"name" bson:"name"`
	Age    int    `json:"age" bson:"age"`
}

type PutReqBody struct {
	Document PutReqBodyDoc `json:"document"`
}

type CollectReqBody struct {
	Collection_name string `json:"collection_name"`
}

type PutRespDocColl struct {
	Collection string        `json:"collection_name" bson:"collection_name"`
	Document   PutReqBodyDoc `json:"document" bson:"document"`
}

type PutRespBody struct {
	Ok bool `json:"ok"`
}
