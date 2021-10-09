package mongodb

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
    "time"
)

var DB *mongo.Database

func InitMongoDb(username, password, host, port, database string)  {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    dbLink := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbLink))
    //defer client.Disconnect(ctx)
    if err != nil {
        log.Fatal("mongo connect error!")
        return
    }
    DB = client.Database(database)
}

func GetMongoDb() *mongo.Database {
    return DB
}

func NewMongoDb()  {

}

type (
    MongoDbOrm struct {

    }
)

func (o *MongoDbOrm) Open() {

}