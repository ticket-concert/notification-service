package handlers_test

import (
	"notification-service/internal/modules/notification/handlers"
	mockcert "notification-service/mocks/modules/notification"
	mocklog "notification-service/mocks/pkg/log"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventHandlerTestSuite struct {
	suite.Suite
	notificationUsecaseCommand *mockcert.UsecaseCommand
	mockLogger                 *mocklog.Logger
}

func (suite *EventHandlerTestSuite) SetupTest() {
	suite.notificationUsecaseCommand = new(mockcert.UsecaseCommand)
	suite.mockLogger = &mocklog.Logger{}
}

func (suite *WorkerHandlerTestSuite) TestInitNotificationEventConflHandler() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailTicket", mock.Anything, mock.Anything).Return(nil)
	suite.notificationUsecaseCommand.On("SendEmailOtpRegister", mock.Anything, mock.Anything).Return(nil)
	handlers.InitNotificationEventConflHandler(suite.notificationUsecaseCommand, suite.mockLogger)
}
