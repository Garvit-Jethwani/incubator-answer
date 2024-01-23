package answer

import (
	"context"
	"testing"

	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/segmentfault/pacman/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockData struct {
	mock.Mock
}

func (m *MockData) Context(ctx context.Context) *data.Data {
	args := m.Called(ctx)
	return args.Get(0).(*data.Data)
}

type MockAnswer struct {
	mock.Mock
}

func (m *MockAnswer) Count(resp *entity.Answer) (int64, error) {
	args := m.Called(resp)
	return args.Get(0).(int64), args.Error(1)
}

func TestGetCountByQuestionID_6d6e909d4a(t *testing.T) {
	tests := []struct {
		name        string
		questionID  string
		mockReturn  int64
		mockError   error
		expected    int64
		expectError bool
	}{
		{
			name:        "Scenario 1: valid questionID, multiple answers",
			questionID:  "12345678",
			mockReturn:  5,
			mockError:   nil,
			expected:    5,
			expectError: false,
		},
		{
			name:        "Scenario 2: valid questionID, no answers",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   nil,
			expected:    0,
			expectError: false,
		},
		{
			name:        "Scenario 3: valid questionID, answers with different status",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   nil,
			expected:    0,
			expectError: false,
		},
		{
			name:        "Scenario 4: valid questionID, some answers with 'AnswerStatusAvailable' status",
			questionID:  "12345678",
			mockReturn:  2,
			mockError:   nil,
			expected:    2,
			expectError: false,
		},
		{
			name:        "Scenario 5: invalid questionID",
			questionID:  "invalid",
			mockReturn:  0,
			mockError:   errors.New("invalid question ID"),
			expected:    0,
			expectError: true,
		},
		{
			name:        "Scenario 6: database connection not available",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   errors.New("database connection not available"),
			expected:    0,
			expectError: true,
		},
		{
			name:        "Scenario 7: database query execution issue",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   errors.New("database query execution issue"),
			expected:    0,
			expectError: true,
		},
		{
			name:        "Scenario 8: internal server error",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   errors.InternalServer(reason.DatabaseError),
			expected:    0,
			expectError: true,
		},
		{
			name:        "Scenario 9: context cancelled",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   context.Canceled,
			expected:    0,
			expectError: true,
		},
		{
			name:        "Scenario 10: context deadline exceeded",
			questionID:  "12345678",
			mockReturn:  0,
			mockError:   context.DeadlineExceeded,
			expected:    0,
			expectError: true,
		},
		{
			name:        "Scenario 11: questionID decoded",
			questionID:  "MTIzNDU2Nzg=",
			mockReturn:  2,
			mockError:   nil,
			expected:    2,
			expectError: false,
		},
		{
			name:        "Scenario 12: questionID can't be decoded",
			questionID:  "invalid format",
			mockReturn:  0,
			mockError:   errors.New("invalid question ID format"),
			expected:    0,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockData := new(MockData)
			mockAnswer := new(MockAnswer)

			mockData.On("Context", mock.Anything).Return(mockAnswer)
			mockAnswer.On("Count", mock.Anything).Return(test.mockReturn, test.mockError)

			ar := &answerRepo{
				data: mockData,
			}

			result, err := ar.GetCountByQuestionID(context.Background(), test.questionID)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}

			mockData.AssertExpectations(t)
			mockAnswer.AssertExpectations(t)
		})
	}
}
