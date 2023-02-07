package gmongo

import (
	"context"
	"proxy-forward/config"
	"proxy-forward/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var db *mongo.Database

// mongodb connect init
func Setup() error {
	//  load options
	opt := options.Client().ApplyURI(config.RuntimeViper.GetString("mongo.uri"))
	// connect mongoDB
	clt, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		logging.Log.Errorf("mongoDB connect Error!  URI: %s", config.RuntimeViper.GetString("mongo.uri"))
		return err
	}
	// Ping MongoDB
	if err := clt.Ping(context.TODO(), nil); err != nil {
		logging.Log.Errorf("mongoDB ping Error!  URI: %s", config.RuntimeViper.GetString("mongo.uri"))
		return err
	}
	// select database
	client = clt
	db = clt.Database(config.RuntimeViper.GetString("mongo.database"))
	return nil
}

// save  {"timestamp": 0, "remote_addr": "", "usage": "", "user_token_id": 1, "type": "req/resp", "uid": 0, "ps_id":  ""}
func SaveForwardData(data bson.M) error {
	// save or not
	if config.RuntimeViper.GetBool("mongo.status") == false {
		return nil
	}
	// Ping MongoDB
	if err := client.Ping(context.TODO(), nil); err != nil {
		logging.Log.Errorf("mongoDB ping Error!")
		return err
	}
	// select collection(proxy_forward)
	col := db.Collection("proxy_forward")
	// insert
	_, err := col.InsertOne(context.TODO(), data)
	if err != nil {
		logging.Log.Warnf("insert forward data failed, ERR: %s", err.Error())
		return err
	}
	return nil
}
