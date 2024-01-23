package answer

import (
	"context"
	"testing"

	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/segmentfault/pacman/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockData struct {
	mock.Mock
}

type MockUniqueIDRepo struct {
	mock.Mock
}

type MockUserRankRepo struct {
	mock.Mock
}

type MockActivityRepo struct {
	mock.Mock
}

func TestUpdateAcceptedStatus_b09402f2e7(t *testing.T) {
	mockData := new(MockData)
	mockUniqueIDRepo := new(MockUniqueIDRepo)
	mockUserRankRepo := new(MockUserRankRepo)
	mockActivityRepo := new(MockActivityRepo)

	ar := &answerRepo{
		data:         mockData,
		uniqueIDRepo: mockUniqueIDRepo,
		userRankRepo: mockUserRankRepo,
		activityRepo: mockActivityRepo,
	}

	testCases := []struct {
		name             string
		ctx              context.Context
		acceptedAnswerID string
		questionID       string
		setup            func()
		shouldError      bool
	}{
		{
			name:             "when context is nil",
			ctx:              nil,
			acceptedAnswerID: "validID",
			questionID:       "validID",
			setup: func() {
				mockData.On("DB.Context", nil).Return(nil, errors.InternalServer(reason.DatabaseError))
			},
			shouldError: true,
		},
		{
			name:             "when acceptedAnswerID and questionID are valid",
			ctx:              context.Background(),
			acceptedAnswerID: "validID",
			questionID:       "validID",
			setup: func() {
				mockData.On("DB.Context", context.Background()).Return(nil, nil)
			},
			shouldError: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := ar.UpdateAcceptedStatus(tt.ctx, tt.acceptedAnswerID, tt.questionID)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
