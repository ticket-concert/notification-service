package usecases_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"notification-service/internal/modules/notification"
	notificationRequest "notification-service/internal/modules/notification/models/request"
	nc "notification-service/internal/modules/notification/usecases"
	"notification-service/internal/modules/order/models/entity"
	"notification-service/internal/pkg/helpers"
	mockcertOrder "notification-service/mocks/modules/order"
	mockhelpers "notification-service/mocks/pkg/helpers"
	mocklog "notification-service/mocks/pkg/log"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockLogger               *mocklog.Logger
	mockEmailSender          *mockhelpers.MailHelper
	mockOrderRepositoryQuery *mockcertOrder.MongodbRepositoryQuery
	usecase                  notification.UsecaseCommand
	ctx                      context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockLogger = &mocklog.Logger{}
	suite.mockEmailSender = &mockhelpers.MailHelper{}
	suite.mockOrderRepositoryQuery = &mockcertOrder.MongodbRepositoryQuery{}
	suite.ctx = context.Background()
	suite.usecase = nc.NewCommandUsecase(
		suite.mockLogger,
		suite.mockEmailSender,
		suite.mockOrderRepositoryQuery,
	)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

func (suite *CommandUsecaseTestSuite) TestSendEmailOtpRegisterSuccess() {
	// Arrange request register
	payload := notificationRequest.SendEmailRegister{
		FullName: "alif septian",
		Email:    "alifseptian@gmail.com",
		UserId:   "1234",
		Otp:      "123456",
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)

	// Act
	err := suite.usecase.SendEmailOtpRegister(suite.ctx, payload)
	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailOtpRegisterFailed() {
	// Arrange request register
	payload := notificationRequest.SendEmailRegister{
		FullName: "alif septian",
		Email:    "alifseptian@gmail.com",
		UserId:   "1234",
		Otp:      "123456",
	}

	expectedError := errors.New("Send email failed")

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(expectedError)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailOtpRegister(suite.ctx, payload)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketSuccess() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrOrder() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
		Error: errors.New("error"),
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrParseOrder() {
	mockOrderById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrNilOrder() {
	mockOrderById := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrEncode() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
		},
		Error: nil,
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrParseFile() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, errors.New("error")
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrGenPdf() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, errors.New("error")
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrWriteFile() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return errors.New("error")
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErrRemoveFile() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return errors.New("error")
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(nil)
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestSendEmailTicketErr() {
	mockOrderById := helpers.Result{
		Data: &entity.Order{
			OrderId:      "id",
			PaymentId:    "id",
			MobileNumber: "08112323123",
			OrderTime:    time.Now(),
			DateTime:     time.Now(),
			Country: entity.Country{
				Place: "place",
			},
			TicketNumber: "1",
		},
	}

	nc.GeneratePDF = func(templatePath string, object interface{}) (*[]byte, error) {
		res := []byte("HTML Pdf")
		return &res, nil
	}

	nc.ParseFiles = func(templatePath string, object interface{}) (*string, error) {
		res := "email"
		return &res, nil
	}

	nc.WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return nil
	}

	nc.Remove = func(name string) error {
		return nil
	}

	suite.mockEmailSender.On("SendEmail", mock.Anything).Return(errors.New("error"))
	suite.mockOrderRepositoryQuery.On("FindOrderById", mock.Anything, mock.Anything).Return(mockChannel(mockOrderById))
	suite.mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything)
	suite.mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything)

	// Act
	err := suite.usecase.SendEmailTicket(suite.ctx, mock.Anything)
	// Assert
	assert.Error(suite.T(), err)
}

// Helper function to create a channel
func mockChannel(result helpers.Result) <-chan helpers.Result {
	responseChan := make(chan helpers.Result)

	go func() {
		responseChan <- result
		close(responseChan)
	}()

	return responseChan
}
