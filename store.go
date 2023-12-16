package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx          = context.Background()
	client       *mongo.Client
	databaseName = "bookstore"
	collection   = "books"
)

type Books struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title"`
	Author string             `json:"author"`
	ISBN   string             `json:"isbn"`
	Price  float64            `json:"price"`
}

func Main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Get("/books/{id}", GetBook)
	r.Post("/books", CreateBook)
	r.Put("/books/{id}", UpdateBook)
	r.Delete("/books/{id}", DeleteBook)

	http.ListenAndServe(":3000", r)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	booksCollection := client.Database(databaseName).Collection(collection)
	var book Book
	err = booksCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	booksCollection := client.Database(databaseName).Collection(collection)
	_, err = booksCollection.InsertOne(ctx, book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	booksCollection := client.Database(databaseName).Collection(collection)
	update := bson.M{
		"$set": bson.M{
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
			"price":  book.Price,
		},
	}
	_, err = booksCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	booksCollection := client.Database(databaseName).Collection(collection)
	_, err = booksCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
