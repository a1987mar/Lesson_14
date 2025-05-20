package main

import (
	"log/slog"
	"net/http"

	"github.com/avpetkun/jessy-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	coll *mongo.Collection
}

type Document struct {
	UserId string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Age    int    `json:"age" bson:"age"`
}

type NewCollect struct {
	Collection_name string `bson:"collection_name"`
}

type DocumentCollection struct {
	CollectionName string     `bson:"collection_name"`
	Document       []Document `bson:"document"`
}

type Doc struct {
	Document Document `bson:"document"`
}

func NewDocumentHandler(c *mongo.Client) *Handler {
	db := c.Database("app")
	coll := db.Collection("collections")
	return &Handler{coll: coll}
}

type PutReqBodyDoc struct {
	UserId string `json:"user_id" bson:"user_id"`
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

func (h *Handler) handlePutDocument(w http.ResponseWriter, r *http.Request) {

	reqBody := PutRespDocColl{}
	err := jessy.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		JSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqBody.Document.Age < 0 || reqBody.Document.Age > 100 {
		JSON(w, "age specified is incorrect", http.StatusBadRequest)
		return
	}

	if len(reqBody.Document.Name) > 18 {
		JSON(w, "name specified is incorrect", http.StatusBadRequest)
		return
	}

	if len(reqBody.Document.UserId) == 0 {
		JSON(w, "user_id is required", http.StatusBadRequest)
		return
	}

	var docMongo DocumentCollection
	filter := bson.M{"collection_name": reqBody.Collection}
	err = h.coll.FindOne(r.Context(), filter).Decode(&docMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			JSON(w, "Document not found", http.StatusNotFound)
		} else {
			slog.Error("Search error:", slog.Any("msg", err.Error()))
			JSON(w, map[string]string{"Search error:": err.Error()}, http.StatusBadRequest)
		}
		return
	}
	for _, v := range docMongo.Document {

		if v.UserId == reqBody.Document.UserId {
			JSON(w, "user with this id EXISTS", http.StatusBadRequest)
			return
		}

	}
	docMongo.Document = append(docMongo.Document, Document(reqBody.Document))
	update := bson.M{"$set": bson.M{"document": docMongo.Document}}
	opts := options.Update().SetUpsert(true)

	_, err = h.coll.UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		JSON(w, map[string]string{"failed to update document:": err.Error()}, http.StatusBadRequest)
		return
	}

	respBody := PutRespBody{Ok: true}
	err = jessy.NewEncoder(w).Encode(respBody)
	if err != nil {
		JSON(w, map[string]string{"failed to encode resp body": err.Error()}, http.StatusOK)
	}
}

func (h *Handler) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	delIdDoc := r.PathValue("idDoc")
	delNameCol := r.PathValue("nameCol")
	if delIdDoc == "" {
		JSON(w, "Document not found", http.StatusBadRequest)
		return
	}
	if delNameCol == "" {
		JSON(w, "Document not found", http.StatusBadRequest)
		return
	}
	var docMongo DocumentCollection
	filter := bson.M{"collection_name": delNameCol}
	err := h.coll.FindOne(r.Context(), filter).Decode(&docMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			JSON(w, "Document not found", http.StatusNotFound)
		} else {
			slog.Error("Search error:", slog.Any("msg", err.Error()))
			JSON(w, map[string]string{"Search error:": err.Error()}, http.StatusBadRequest)
		}
		return
	}
	initLenDocMongo := len(docMongo.Document)
	docCopy := []Document{}
	for _, v := range docMongo.Document {
		if v.UserId != delIdDoc {
			docCopy = append(docCopy, v)
		}
	}
	if initLenDocMongo == len(docCopy) {
		JSON(w, "user not found", http.StatusNotFound)
		return
	}
	docMongo.Document = docCopy
	update := bson.M{"$set": docMongo}
	opts := options.Update().SetUpsert(true)
	_, err = h.coll.UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		JSON(w, map[string]string{"failed to update document:": err.Error()}, http.StatusBadRequest)
		return
	}

	respBody := PutRespBody{Ok: true}
	err = jessy.NewEncoder(w).Encode(respBody)
	if err != nil {
		JSON(w, map[string]string{"failed to encode resp body": err.Error()}, http.StatusOK)
	}
}

func JSON(w http.ResponseWriter, data any, statuscode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	err := jessy.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("JSON encode error:", slog.Any("msg", err.Error()))
		return
	}
}

func (h *Handler) handleNewCollect(w http.ResponseWriter, r *http.Request) {
	reqBody := CollectReqBody{}
	err := jessy.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		JSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	clt := &NewCollect{
		Collection_name: reqBody.Collection_name,
	}
	filter := bson.M{"collection_name": clt.Collection_name}
	update := bson.M{"$set": clt}
	opts := options.Update().SetUpsert(true)

	_, err = h.coll.UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		JSON(w, map[string]string{"failed to update collection:": err.Error()}, http.StatusBadRequest)
		return
	}
	respBody := PutRespBody{Ok: true}
	err = jessy.NewEncoder(w).Encode(respBody)
	if err != nil {
		JSON(w, map[string]string{"failed to encode resp body": err.Error()}, http.StatusOK)
		return
	}
}

func (h *Handler) handleListCollections(w http.ResponseWriter, r *http.Request) {
	cursor, err := h.coll.Find(r.Context(), bson.M{})
	if err != nil {
		JSON(w, map[string]string{"failed to fetch documents:": err.Error()}, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())
	var coll []NewCollect
	if err := cursor.All(r.Context(), &coll); err != nil {
		JSON(w, map[string]string{"failed to decode documents: ": err.Error()}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	JSON(w, coll, http.StatusOK)
}

func (h *Handler) handleGetCollect(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		JSON(w, "collection name is required", http.StatusBadRequest)
		return
	}
	filter := bson.M{"collection_name": name}
	var coll DocumentCollection
	err := h.coll.FindOne(r.Context(), filter).Decode(&coll)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			JSON(w, "Document not found", http.StatusNotFound)
		} else {
			slog.Error("Search error:", slog.Any("msg", err.Error()))
			JSON(w, map[string]string{"Search error:": err.Error()}, http.StatusBadRequest)
		}
		return
	}
	JSON(w, coll, http.StatusOK)
}

func (h *Handler) handleDeleteCollect(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		JSON(w, "collection name is required", http.StatusBadRequest)
		return
	}
	filter := bson.M{"collection_name": name}
	result, err := h.coll.DeleteOne(r.Context(), filter)
	if err != nil {
		slog.Info("DeleteOne error:", slog.Any("msg", err.Error()))
		JSON(w, err.Error(), http.StatusOK)
		return
	}
	if result.DeletedCount == 0 {
		JSON(w, "collection not found", http.StatusNotFound)
		return
	}
	JSON(w, "collection delete", http.StatusOK)

}
