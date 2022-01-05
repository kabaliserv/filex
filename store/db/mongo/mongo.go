package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func test() {
	client, err := mongo.NewClient(options.Client().ApplyURI("<ATLAS_URI_HERE>"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	coll := client.Database("tea").Collection("ratings")
	docs := []interface{}{
		bson.D{{"type", "Masala"}, {"rating", 10}},
		bson.D{{"type", "Matcha"}, {"rating", 7}},
		bson.D{{"type", "Assam"}, {"rating", 4}},
		bson.D{{"type", "Oolong"}, {"rating", 9}},
		bson.D{{"type", "Chrysanthemum"}, {"rating", 5}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}},
		bson.D{{"type", "Jasmine"}, {"rating", 3}},
		bson.D{{"type", "English Breakfast"}, {"rating", 6}},
		bson.D{{"type", "White Peony"}, {"rating", 4}},
	}
	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents inserted: %d\n", len(result.InsertedIDs))
}
