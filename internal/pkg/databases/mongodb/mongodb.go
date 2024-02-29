package mongodb

import (
	"context"
	"notification-service/configs"
	"strconv"
	"strings"

	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

var mongoMasterClient *mongo.Client
var mongoMasterDbName string
var mongoSlaveClient *mongo.Client
var mongoSlaveDbName string

type MongoCollections interface {
	InitConnection(masterDBUrl, slaveDBUrl string) (*mongo.Client, *mongo.Client)
	NewClient(mongoUri string) (*mongo.Client, error)
	NewClientSlave(mongoUri string) (*mongo.Client, error)
}

type MongoImpl struct {
	Collections MongoCollections
}

func (m *MongoImpl) InitConnection(masterDBUrl, slaveDBUrl string) (masterClient, slaveClient *mongo.Client) {
	mClient, err := m.Collections.NewClient(masterDBUrl)

	if err != nil {
		panic(err)
	}

	sClient, err := m.Collections.NewClientSlave(slaveDBUrl)

	if err != nil {
		panic(err)
	}

	mongoMasterClient = mClient
	mongoMasterDbName = getDbName(masterDBUrl)
	mongoSlaveClient = sClient
	mongoSlaveDbName = getDbName(slaveDBUrl)

	return mClient, sClient
}

func (m *MongoImpl) SetCollections(collections MongoCollections) {
	m.Collections = collections
}

func GetMasterConn() *mongo.Client {
	return mongoMasterClient
}

func GetMasterDBName() string {
	return mongoMasterDbName
}

func GetSlaveConn() *mongo.Client {
	return mongoSlaveClient
}

func GetSlaveDBName() string {
	return mongoSlaveDbName
}

func getDbName(s string) string {
	ss := strings.Split(s, "?")
	if len(ss) > 1 {
		return strings.Split(strings.ReplaceAll(ss[0], "//", ""), "/")[1]
	} else {
		return strings.Split(strings.ReplaceAll(s, "//", ""), "/")[1]
	}
}

func (m *MongoImpl) NewClient(mongoUri string) (*mongo.Client, error) {
	var monitor *event.CommandMonitor
	if configs.GetConfig().Datadog.DatadogEnabled == "true" {
		monitor = mongotrace.NewMonitor()
	} else {
		monitor = apmmongo.CommandMonitor()
	}
	var mongoPoolSize uint64
	mongoPoolSize = 100
	parseRedisDb, err := strconv.ParseInt(configs.GetConfig().MongoDB.MongoPoolSize, 10, 32)

	if err == nil {
		mongoPoolSize = uint64(parseRedisDb)
	}
	client, err := mongo.Connect(
		context.Background(),
		options.Client().SetMonitor(monitor),
		options.Client().SetRetryWrites(true),
		options.Client().SetRetryReads(true),
		options.Client().SetMaxPoolSize(mongoPoolSize),
		options.Client().SetMinPoolSize(20),
		options.Client().ApplyURI(mongoUri),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *MongoImpl) NewClientSlave(mongoUri string) (*mongo.Client, error) {
	var monitor *event.CommandMonitor
	if configs.GetConfig().Datadog.DatadogEnabled == "true" {
		monitor = mongotrace.NewMonitor()
	} else {
		monitor = apmmongo.CommandMonitor()
	}

	client, err := mongo.Connect(
		context.Background(),
		options.Client().SetMonitor(monitor),
		options.Client().SetRetryWrites(true),
		options.Client().SetRetryReads(true),
		options.Client().SetMaxPoolSize(100),
		options.Client().SetMinPoolSize(50),
		options.Client().ApplyURI(mongoUri),
		options.Client().SetReadPreference(readpref.SecondaryPreferred()),
	)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), readpref.SecondaryPreferred()); err != nil {
		return nil, err
	}

	return client, nil
}
