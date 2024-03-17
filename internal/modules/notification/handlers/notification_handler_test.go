package handlers_test

import (
	"notification-service/internal/modules/notification/handlers"
	"notification-service/internal/pkg/errors"
	mockcert "notification-service/mocks/modules/notification"
	mocklog "notification-service/mocks/pkg/log"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type WorkerHandlerTestSuite struct {
	suite.Suite
	notificationUsecaseCommand *mockcert.UsecaseCommand
	handler                    *handlers.NotificationEventHandler
	mockLogger                 *mocklog.Logger
}

func (suite *WorkerHandlerTestSuite) SetupTest() {
	suite.notificationUsecaseCommand = new(mockcert.UsecaseCommand)
	suite.mockLogger = &mocklog.Logger{}
	suite.handler = &handlers.NotificationEventHandler{
		suite.notificationUsecaseCommand,
		suite.mockLogger,
	}
	handlers.NewInsertNotificationConsumer(suite.notificationUsecaseCommand, suite.mockLogger)
}

func TestWorkerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(WorkerHandlerTestSuite))
}

func (suite *WorkerHandlerTestSuite) TestSendEmailOtpRegister() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailOtpRegister", mock.Anything, mock.Anything).Return(nil)
	topic := "topic"
	msg := kafka.Message{
		Value: []byte(`{"userId": "154fb7e0-07b8-43ba-898c-2de889f272ef",
		"fullName": "fullName", "email": "alif@gmail.com", "otp":"123"}`),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    kafka.OffsetBeginning,
		},
	}
	suite.handler.SendEmailOtpRegister(&msg, topic)
}

func (suite *WorkerHandlerTestSuite) TestSendEmailOtpRegisterErrParse() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailOtpRegister", mock.Anything, mock.Anything).Return(nil)
	topic := "topic"
	msg := kafka.Message{
		Value: []byte("test"),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    kafka.OffsetBeginning,
		},
	}
	suite.handler.SendEmailOtpRegister(&msg, topic)
}

func (suite *WorkerHandlerTestSuite) TestSendEmailOtpRegisterErr() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailOtpRegister", mock.Anything, mock.Anything).Return(errors.BadRequest("error"))
	topic := "topic"
	msg := kafka.Message{
		Value: []byte(`{"userId": "154fb7e0-07b8-43ba-898c-2de889f272ef",
		"fullName": "fullName", "email": "alif@gmail.com", "otp":"123"}`),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    kafka.OffsetBeginning,
		},
	}
	suite.handler.SendEmailOtpRegister(&msg, topic)
}

func (suite *WorkerHandlerTestSuite) TestSendEmailTicket() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailTicket", mock.Anything, mock.Anything).Return(nil)
	topic := "topic"
	handlers.Unmarshal = func(data []byte, v any) error {
		return nil
	}
	msg := kafka.Message{
		Value: []byte("154fb7e0-07b8-43ba-898c-2de889f272ef"),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    kafka.OffsetBeginning,
		},
	}
	suite.handler.SendEmailTicket(&msg, topic)
}

func (suite *WorkerHandlerTestSuite) TestSendEmailTicketErrMsg() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailTicket", mock.Anything, mock.Anything).Return(nil)
	topic := "topic"
	handlers.Unmarshal = func(data []byte, v any) error {
		return errors.BadRequest("error")
	}
	msg := kafka.Message{
		Value: []byte(""),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    kafka.OffsetBeginning,
		},
	}
	suite.handler.SendEmailTicket(&msg, topic)
}

func (suite *WorkerHandlerTestSuite) TestSendEmailTicketErr() {
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)
	suite.notificationUsecaseCommand.On("SendEmailTicket", mock.Anything, mock.Anything).Return(errors.BadRequest("error"))
	topic := "topic"
	handlers.Unmarshal = func(data []byte, v any) error {
		return nil
	}
	msg := kafka.Message{
		Value: []byte("154fb7e0-07b8-43ba-898c-2de889f272ef"),
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 1,
			Offset:    kafka.OffsetBeginning,
		},
	}
	suite.handler.SendEmailTicket(&msg, topic)
}
