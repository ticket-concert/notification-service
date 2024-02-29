package apm_test

import (
	"notification-service/internal/pkg/apm"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ApmTestSuite struct {
	suite.Suite
}

func (suite *ApmTestSuite) SetupTest() {
	os.Unsetenv("ELASTIC_APM_SERVER_URL")
	os.Unsetenv("ELASTIC_APM_SECRET_TOKEN")
}

func TestApmTestSuite(t *testing.T) {
	suite.Run(t, new(ApmTestSuite))
}

func (suite *ApmTestSuite) TestInitConnection() {
	// Set up test data
	apmUrl := "http://example.com"
	apmSecretToken := "secretToken123"

	// Call the function under test
	apm.InitConnection(apmUrl, apmSecretToken)

	// Assertions
	assert.Equal(suite.T(), apmUrl, os.Getenv("ELASTIC_APM_SERVER_URL"))
	assert.Equal(suite.T(), apmSecretToken, os.Getenv("ELASTIC_APM_SECRET_TOKEN"))
}

func (suite *ApmTestSuite) TestGetTracer() {
	// Set up test data
	os.Setenv("SERVICE_NAME", "user_service")
	os.Setenv("SERVICE_VERSION", "1.0.0")

	// Call the function under test
	tracer := apm.GetTracer()

	// Assertions
	assert.NotNil(suite.T(), tracer)
}
