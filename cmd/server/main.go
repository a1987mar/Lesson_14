package main

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	opts := options.Client()
	opts.ApplyURI("mongodb://root:root@localhost:27017")
	clt, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	defer clt.Disconnect(ctx)
	handlers := NewDocumentHandler(clt)

	http.HandleFunc("POST /put_document", handlers.handlePutDocument)
	http.HandleFunc("DELETE /delete_document/{nameCol}/{idDoc}", handlers.handleDeleteDocument)

	http.HandleFunc("POST /create_collection", handlers.handleNewCollect)
	http.HandleFunc("GET /get_collection/{name}", handlers.handleGetCollect)
	http.HandleFunc("GET /list_collections", handlers.handleListCollections)
	http.HandleFunc("DELETE /delete_collection/{name}", handlers.handleDeleteCollect)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(fmt.Errorf("server listening failed: %v", err))
	}
}
