package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"notification-service/internal/modules/notification"
	notificationRequest "notification-service/internal/modules/notification/models/request"
	nc "notification-service/internal/modules/notification/usecases"
	mockhelpers "notification-service/mocks/pkg/helpers"
	mocklog "notification-service/mocks/pkg/log"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockLogger      *mocklog.Logger
	mockEmailSender *mockhelpers.EmailSender
	usecase         notification.UsecaseCommand
	ctx             context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockLogger = &mocklog.Logger{}
	suite.mockEmailSender = &mockhelpers.EmailSender{}
	suite.ctx = context.WithValue(context.TODO(), "key", "value")
	suite.usecase = nc.NewCommandUsecase(
		suite.mockLogger,
		suite.mockEmailSender,
	)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

func (suite *CommandUsecaseTestSuite) TestSendEmailOtpRegisterSuccess() {
	// Arrange request register
	payload := notificationRequest.SendEmailRegister{
		FullName: "Irman Juliansyah",
		Email:    "irmanjuliansyah@gmail.com",
		UserId:   "1234",
		Otp:      "123456",
	}

	suite.mockEmailSender.On("SendEmailSimple", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	err := suite.usecase.SendEmailOtpRegister(suite.ctx, payload)
	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailOtpRegisterFailed() {
	// Arrange request register
	payload := notificationRequest.SendEmailRegister{
		FullName: "Irman Juliansyah",
		Email:    "irmanjuliansyah@gmail.com",
		UserId:   "1234",
		Otp:      "123456",
	}

	expectedError := errors.New("Send email failed")

	suite.mockEmailSender.On("SendEmailSimple", mock.Anything, mock.Anything, mock.Anything).Return(expectedError)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailOtpRegister(suite.ctx, payload)
	// Assert
	assert.EqualError(suite.T(), err, expectedError.Error())
}
