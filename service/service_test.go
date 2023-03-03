package service

// import (
// 	"context"
// 	"testing"

// 	db "github.com/amancooks08/BookMySport/db"
// 	"github.com/amancooks08/BookMySport/mocks"
// 	"github.com/stretchr/testify/suite"
// )

// type ServiceTestSuite struct {
// 	suite.Suite
// 	mockService Services
// 	repository  *mocks.Storer
// }

// func TestServiceTestSuite(t *testing.T) {
// 	suite.Run(t, new(ServiceTestSuite))
// }

// func (suite *ServiceTestSuite) SetupTest() {
// 	suite.repository = new(mocks.Storer)
// 	suite.mockService = NewCustomerOps(suite.repository)
// }

// func (suite *ServiceTestSuite) TearDownTest() {
// 	suite.repository.AssertExpectations(suite.T())
// }

// func (suite *ServiceTestSuite) TestRegisterUser() {
// 	t := suite.T()

// 	type args struct {
// 		ctx  context.Context
// 		user *db.User
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 		prepare func(args, *mocks.Storer)
// 	}{
// 		{
// 			name: "",
// 			args: args{
// 				ctx: context.TODO(),
// 				user: &db.User{
// 					Name : "Amandeep Singh",
// 					Contact: "987654321",
// 					Email:   "amandeep.singh2001@gmail.com",
// 					Password: "123456",
// 					City: "Delhi",
// 					State: "Delhi",
// 					Role: "customer",
// 				},
// 			},
// 			wantErr: false,
// 			prepare: func(a args, m *mocks.Storer) {
// 				m.On("RegisterCustomer", a.context.TODO(), mock.Anything).Return(nil).Once()
// 			}
// 		},
// 		{
// 			name: "Register User",
// 			args: args{
// 				ctx: context.Background(),
// 				user: &db.User{
// 					Name : "Amandeep Singh Chhabra",
// 					Contact: "9876543210",
// 					Email:   "test1@gmail.com",
// 					Password: "123456",
// 					City: "Delhi",
// 					State: "Delhi",
// 					Role: "customer",
// 				},
// 			},

// 		}

// }
// func TestUserOps_LoginUser(t *testing.T) {

// }
