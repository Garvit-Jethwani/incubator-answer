package answer

import (
	"context"
	"errors"
	"testing"

	"github.com/apache/incubator-answer/internal/entity"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockData struct {
	mock.Mock
}

func (m *MockData) Context(ctx context.Context) *xorm.Session {
	args := m.Called(ctx)
	return args.Get(0).(*xorm.Session)
}

type MockSession struct {
	mock.Mock
}

func (m *MockSession) ID(id interface{}) *xorm.Session {
	args := m.Called(id)
	return args.Get(0).(*xorm.Session)
}

func (m *MockSession) Cols(columns ...string) *xorm.Session {
	args := m.Called(columns)
	return args.Get(0).(*xorm.Session)
}

func (m *MockSession) Update(bean interface{}, condiBean ...interface{}) (int64, error) {
	args := m.Called(bean, condiBean)
	return args.Get(0).(int64), args.Error(1)
}

type MockRepo struct {
	mock.Mock
	data *MockData
}

func (m *MockRepo) updateSearch(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUpdateAnswer(t *testing.T) {
	tests := []struct {
		name          string
		answer        *entity.Answer
		cols          []string
		mockDB        func(m *MockData, s *MockSession)
		mockRepo      func(m *MockRepo)
		expectedError error
	}{
		{
			name: "Test scenario 1: Successful update",
			answer: &entity.Answer{
				ID:         "1",
				QuestionID: "1",
				Status:     1,
			},
			cols: []string{"status"},
			mockDB: func(m *MockData, s *MockSession) {
				m.On("Context", mock.Anything).Return(s)
				s.On("ID", mock.Anything).Return(s)
				s.On("Cols", mock.Anything).Return(s)
				s.On("Update", mock.Anything, mock.Anything).Return(1, nil)
			},
			mockRepo: func(m *MockRepo) {
				m.On("updateSearch", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		// TODO: Add more test cases here for other scenarios following the same pattern.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockData := new(MockData)
			mockSession := new(MockSession)
			mockRepo := new(MockRepo)
			mockRepo.data = mockData

			tt.mockDB(mockData, mockSession)
			tt.mockRepo(mockRepo)

			ar := &answerRepo{
				data: mockRepo.data,
			}
			err := ar.UpdateAnswer(context.Background(), tt.answer, tt.cols)

			assert.Equal(t, tt.expectedError, err)
			mockData.AssertExpectations(t)
			mockSession.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
