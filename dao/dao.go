package dao

import (
	"context"
	"fmt"
	"log"

	"github.com/LuisCusihuaman/phonebook-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CONNECTIONSTRING DB connection string
const CONNECTIONSTRING = "MONGODBURI"

// DBNAME Database name
const DBNAME = "phonebook"

// COLLNAME Collection name
const COLLNAME = "people"

var db *mongo.Database

// Connect establish a connection to database
func init() {
	clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	db = client.Database(DBNAME)

	fmt.Println("Connected to MongoDB!")
	// Collection types can be used to access the database
}

// InsertManyValues inserts many items from byte slice
func InsertManyValues(people []models.Person) {
	var ppl []interface{}
	for _, p := range people {
		ppl = append(ppl, p)
	}
	_, err := db.Collection(COLLNAME).InsertMany(context.Background(), ppl)
	if err != nil {
		log.Fatal(err)
	}
}

// InsertOneValue inserts one item from Person model
func InsertOneValue(person models.Person) {
	fmt.Println(person)
	_, err := db.Collection(COLLNAME).InsertOne(context.Background(), person)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAllPeople returns all people from DB
func GetAllPeople() []models.Person {
	cur, err := db.Collection(COLLNAME).Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("exploto")
		log.Fatal(err)
	}
	var elements []models.Person
	var elem models.Person
	// Get the next result from the cursor
	for cur.Next(context.Background()) {
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		elements = append(elements, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return elements
}

// DeletePerson deletes an existing person
func DeletePerson(person models.Person) {
	_, err := db.Collection(COLLNAME).DeleteOne(context.Background(), person, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// UpdatePerson updates an existing person
func UpdatePerson(person models.Person, personID string) {
	doc := db.Collection(COLLNAME).FindOneAndUpdate(
		context.Background(), bson.M{
			"id": personID,
			"$set": bson.M{
				"firstname":           person.Firstname,
				"lastname":            person.Lastname,
				"contactinfo.city":    person.City,
				"contactinfo.zipcode": person.Zipcode,
				"contactinfo.phone":   person.Phone}}, nil)
	fmt.Println(doc)
}
