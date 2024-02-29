package mongodb_test

import (
	"notification-service/internal/pkg/databases/mongodb"
	mockmongo "notification-service/mocks/pkg/databases/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSuite struct {
	suite.Suite
	mock *mockmongo.MongoCollections
	impl *mongodb.MongoImpl
}

func (suite *MongoSuite) SetupTest() {
	suite.mock = &mockmongo.MongoCollections{}
	suite.impl = &mongodb.MongoImpl{}
	suite.impl.SetCollections(suite.mock)
}

func (suite *MongoSuite) TestInitConnection() {
	// Set mock
	suite.mock.On("NewClient", mock.Anything).Return(&mongo.Client{}, nil)
	suite.mock.On("NewClientSlave", mock.Anything).Return(&mongo.Client{}, nil)

	// Call the function being tested
	mClient, sClient := suite.impl.InitConnection("mongodb://test:27020/admin", "mongodb://test:27020/admin")

	// Assert
	assert.NotNil(suite.T(), mClient, "expected nil because of valid uri")
	assert.NotNil(suite.T(), sClient, "expected nil because of valid uri")
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
