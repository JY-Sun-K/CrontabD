package main

import (
	"context"
	"fmt"
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

type FindByJobName struct {
	JobName string `bson:"jobName"`
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

	cond := &FindByJobName{JobName: "job10"}

	cursor, err := collection.Find(context.TODO(), cond, options.Find().SetSkip(0), options.Find().SetLimit(2))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		record := &LogRecord{}
		err := cursor.Decode(record)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*record)
	}

}
