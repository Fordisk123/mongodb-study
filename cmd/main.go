package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"time"
)

type Person struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Age        int                `bson:"age,omitempty"`
	Born       time.Time          `bson:"born,omitempty"`
	Hobby      []string           `bson:"hobby,omitempty"`
	Transcript `bson:"transcript,omitempty"`
}

type Transcript struct {
	Math    float32 `bson:"math,omitempty"`
	English float32 `bson:"english,omitempty"`
	Chinese float32 `bson:"chinese,omitempty"`
}

func main() {

	person := Person{
		Name:  "xiaoming",
		Age:   18,
		Born:  time.Now(),
		Hobby: []string{"game", "ball"},
		Transcript: Transcript{
			Math:    90,
			English: 60,
			Chinese: 70,
		},
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:Chison2019%21@192.168.55.116/test?authSource=admin"))
	if err != nil {
		fmt.Println(err)
		return
	}

	collection := client.Database("test").Collection("person")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	var id primitive.ObjectID

	for i := 0; i < 10; i++ {
		res, err := collection.InsertOne(ctx, person)
		if err != nil {
			fmt.Println(err)
			return
		}
		id = res.InsertedID.(primitive.ObjectID)
	}

	fmt.Printf("Res id : %s\n", id.Hex())

	//select by id
	fmt.Println("------------------------ select one by id ------------------------")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{"_id": id}
	err = collection.FindOne(ctx, filter).Decode(&person)
	if err != nil {
		panic(err)
	}
	fmt.Println(person)

	//update
	fmt.Println("------------------------ update ------------------------")
	ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	person.Name = "xiaohong"
	person.Chinese = 100
	person.Hobby = []string{"read", "dance"}
	updateRes, err := collection.UpdateMany(ctx, filter, bson.M{"$set": person})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated records count : %d \n", updateRes.ModifiedCount)

	//select all
	fmt.Println("------------------------ select all ------------------------")
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	var personList []Person
	allRes, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = allRes.All(ctx, &personList)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, person := range personList {
		fmt.Println(person)
	}

	//aggregate
	fmt.Println("------------------------ aggregate ------------------------")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	matchPipe := bson.D{
		{"$match", bson.M{
			"transcript.chinese": bson.M{
				"$gte": 90,
			},
		},
		},
	}
	showInfoCursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchPipe})
	if err != nil {
		panic(err)
	}
	if err = showInfoCursor.All(ctx, &personList); err != nil {
		panic(err)
	}
	for _, person := range personList {
		fmt.Println(person)
	}

	//delete
	fmt.Println("------------------------ delete ------------------------")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter = bson.M{"age": 18}
	delRes, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted records count : %d \n", delRes.DeletedCount)
}
