package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.101.138:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}
	collection := client.Database("cron").Collection("log")

	record := &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{time.Now().Unix(), time.Now().Unix() + 10},
	}

	result, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		fmt.Println(err)
		return
	}

	docId := result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID", docId.Hex())

}
