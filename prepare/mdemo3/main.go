package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// startTime小于某时间
// {"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime": {"$lt": timestamp} }
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
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

	// 4, 要删除开始时间早于当前时间的所有日志($lt是less than)
	//  delete({"timePoint.startTime": {"$lt": 当前时间}})
	delCond := &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	// 执行删除
	delResult, err := collection.DeleteMany(context.TODO(), delCond)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除的行数:", delResult.DeletedCount)
}
