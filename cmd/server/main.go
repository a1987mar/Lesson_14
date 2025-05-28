package main

import (
	"context"
	"fmt"
	"lesson_14/cmd/server/handler"

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
	handlers := handler.NewDocumentHandler(clt)

	http.HandleFunc("POST /put_document", handlers.HandlePutDocument)
	http.HandleFunc("DELETE /delete_document/{nameCol}/{idDoc}", handlers.HandleDeleteDocument)

	http.HandleFunc("POST /create_collection", handlers.HandleNewCollect)
	http.HandleFunc("GET /get_collection/{name}", handlers.HandleGetCollect)
	http.HandleFunc("GET /list_collections", handlers.HandleListCollections)
	http.HandleFunc("DELETE /delete_collection/{name}", handlers.HandleDeleteCollect)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(fmt.Errorf("server listening failed: %v", err))
	}
}
