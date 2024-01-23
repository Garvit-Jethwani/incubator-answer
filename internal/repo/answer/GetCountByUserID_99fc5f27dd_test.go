package answer

import (
	"context"
	"testing"

	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/segmentfault/pacman/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockData struct {
	mock.Mock
}

func (m *MockData) Context(ctx context.Context) *data.Context {
	args := m.Called(ctx)
	return args.Get(0).(*data.Context)
}

type MockContext struct {
	mock.Mock
}

func (m *MockContext) Where(query interface{}, args ...interface{}) *data.Context {
	m.Called(query, args)
	return m
}

func (m *MockContext) Count(bean interface{}) (int64, error) {
	args := m.Called(bean)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetCountByUserID_99fc5f27dd(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		mockData    func(mockData *MockData, mockContext *MockContext)
		expectedErr error
		expectedRes int64
	}{
		{
			name:   "Valid Scenario",
			userID: uid.NewID(),
			mockData: func(mockData *MockData, mockContext *MockContext) {
				mockData.On("Context", mock.Anything).Return(mockContext)
				mockContext.On("Where", mock.Anything, mock.Anything).Return(mockContext)
				mockContext.On("Count", mock.Anything).Return(int64(10), nil)
			},
			expectedErr: nil,
			expectedRes: 10,
		},
		{
			name:   "Invalid Scenario",
			userID: "invalidUserID",
			mockData: func(mockData *MockData, mockContext *MockContext) {
				mockData.On("Context", mock.Anything).Return(mockContext)
				mockContext.On("Where", mock.Anything, mock.Anything).Return(mockContext)
				mockContext.On("Count", mock.Anything).Return(int64(0), nil)
			},
			expectedErr: nil,
			expectedRes: 0,
		},
		{
			name:   "Error Scenario",
			userID: uid.NewID(),
			mockData: func(mockData *MockData, mockContext *MockContext) {
				mockData.On("Context", mock.Anything).Return(mockContext)
				mockContext.On("Where", mock.Anything, mock.Anything).Return(mockContext)
				mockContext.On("Count", mock.Anything).Return(int64(0), errors.InternalServer(reason.DatabaseError))
			},
			expectedErr: errors.InternalServer(reason.DatabaseError),
			expectedRes: 0,
		},
		// TODO: Add more test cases for the remaining scenarios.
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockData := new(MockData)
			mockContext := new(MockContext)
			test.mockData(mockData, mockContext)

			ar := &answerRepo{
				data: mockData,
			}

			res, err := ar.GetCountByUserID(context.Background(), test.userID)
			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, test.expectedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes, res)
			}

			mockData.AssertExpectations(t)
			mockContext.AssertExpectations(t)
		})
	}
}
