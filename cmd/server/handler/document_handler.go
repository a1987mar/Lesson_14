package handler

import (
	"log/slog"
	"net/http"

	"github.com/avpetkun/jessy-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handler) HandlePutDocument(w http.ResponseWriter, r *http.Request) {

	reqBody := PutRespDocColl{}
	err := jessy.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		JSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info(reqBody.Document.Name)
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

	if isInvalidCollectionName(reqBody.Collection) {
		JSON(w, "invalid collection name", http.StatusBadRequest)
		return
	}

	filter := bson.M{"userid": reqBody.Document.UserId}
	update := bson.M{"$set": reqBody.Document}
	opts := options.Update().SetUpsert(true)

	_, err = h.db.Collection(reqBody.Collection).UpdateOne(r.Context(), filter, update, opts)
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

func (h *Handler) HandleDeleteDocument(w http.ResponseWriter, r *http.Request) {
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
	filter := bson.M{"userid": delIdDoc}
	result, err := h.db.Collection(delNameCol).DeleteOne(r.Context(), filter)
	if err != nil {
		JSON(w, "Document not found", http.StatusNotFound)
		return
	}
	if result.DeletedCount == 0 {
		JSON(w, "Document not found", http.StatusNotFound)
		return
	}
	JSON(w, result, http.StatusOK)
}

func (h *Handler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	getNameUser := r.PathValue("nameuser")
	getNameCol := r.PathValue("nameCol")
	if getNameUser == "" || getNameCol == "" {
		JSON(w, "Collection name and user name are required", http.StatusBadRequest)
		return
	}
	filter := bson.M{"name": getNameUser}
	result, err := h.db.Collection(getNameCol).Find(r.Context(), filter)
	if err != nil {
		JSON(w, "Database query failed", http.StatusOK)
	}
	defer result.Close(r.Context())

	var results []Document
	if err := result.All(r.Context(), &results); err != nil {
		JSON(w, map[string]string{"Decode error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if len(results) == 0 {
		JSON(w, "No documents found with that name", http.StatusNotFound)
		return
	}
	JSON(w, results, http.StatusOK)
}
