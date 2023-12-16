package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri     = "mongodb://localhost:27017"
	dbName  = "booksdb"
	colName = "books"
)

type Book struct {
	ID     string  `bson:"_id,omitempty"`
	Title  string  `bson:"title"`
	Author string  `bson:"author"`
	ISBN   string  `bson:"isbn"`
	Price  float64 `bson:"price"`
}

func main() {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(dbName).Collection(colName)

	// Sample book data
	book := Book{
		Title:  "Sample Book",
		Author: "John Doe",
		ISBN:   "1234567890",
		Price:  24.99,
	}

	// Marshal the book object to BSON
	bsonData, err := bson.Marshal(book)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal BSON data back to a book object
	var decodedBook Book
	if err := bson.Unmarshal(bsonData, &decodedBook); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encoded BSON:", bsonData)
	fmt.Println("Decoded Book:", decodedBook)
}
