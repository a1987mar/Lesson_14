package handler

import (
	"log/slog"
	"net/http"

	"github.com/avpetkun/jessy-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) HandleNewCollect(w http.ResponseWriter, r *http.Request) {
	reqBody := CollectReqBody{}
	err := jessy.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		JSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	if reqBody.Collection_name == "" {
		JSON(w, "collection name is required", http.StatusBadRequest)
		return
	}
	if isInvalidCollectionName(reqBody.Collection_name) {
		JSON(w, "invalid collection name", http.StatusBadRequest)
		return
	}
	err = h.db.CreateCollection(r.Context(), reqBody.Collection_name)
	if err != nil {
		JSON(w, map[string]string{"failed to create collection": err.Error()}, http.StatusBadRequest)
		return
	}
	respBody := PutRespBody{Ok: true}
	JSON(w, respBody, http.StatusOK)
}

func (h *Handler) HandleListCollections(w http.ResponseWriter, r *http.Request) {

	collections, err := h.db.ListCollectionNames(r.Context(), bson.M{})
	if err != nil {
		JSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	JSON(w, collections, http.StatusOK)
}

func (h *Handler) HandleGetCollect(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		JSON(w, "collection name is required", http.StatusBadRequest)
		return
	}

	cursor, err := h.db.Collection(name).Find(r.Context(), bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			JSON(w, "Document not found", http.StatusNotFound)
		} else {
			slog.Error("Search error:", slog.Any("msg", err.Error()))
			JSON(w, map[string]string{"Search error:": err.Error()}, http.StatusBadRequest)
		}
		return
	}
	defer cursor.Close(r.Context())
	var docs []Document
	if err = cursor.All(r.Context(), &docs); err != nil {
		JSON(w, map[string]string{"failed to decode documents: ": err.Error()}, http.StatusInternalServerError)
		return
	}
	JSON(w, docs, http.StatusOK)
}

func (h *Handler) HandleDeleteCollect(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		JSON(w, "collection name is required", http.StatusBadRequest)
		return
	}
	err := h.db.Collection(name).Drop(r.Context())
	if err != nil {
		JSON(w, "collection not found", 404)
		return
	}
	JSON(w, "collection delete", http.StatusOK)
}
