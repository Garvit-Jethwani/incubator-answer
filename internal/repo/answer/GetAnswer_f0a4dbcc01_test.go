package answer

import (
	"context"
	"errors"
	"testing"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/base/reason"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/go-xorm/xorm"
	"github.com/segmentfault/pacman/errors"
	"github.com/stretchr/testify/assert"
)

type answerRepo struct {
	data *data.Data
}

func TestGetAnswer_f0a4dbcc01(t *testing.T) {
	ar := &answerRepo{
		data: &data.Data{},
	}

	testCases := []struct {
		name           string
		id             string
		dbMockFunc     func() error
		shortIdEnabled bool
		expectedError  error
		expectedExist  bool
	}{
		{
			name:          "Positive Test Scenario",
			id:            "123",
			dbMockFunc:    func() error { return nil },
			expectedExist: true,
		},
		{
			name:          "Negative Test Scenario",
			id:            "invalid",
			dbMockFunc:    func() error { return nil },
			expectedError: errors.New("invalid id"),
		},
		{
			name:          "Positive Test Scenario with deshortened id",
			id:            uid.EnShortID("123"),
			dbMockFunc:    func() error { return nil },
			expectedExist: true,
		},
		{
			name:          "Negative Test Scenario with non-existing id",
			id:            "999",
			dbMockFunc:    func() error { return nil },
			expectedError: errors.New("answer does not exist"),
		},
		{
			name:           "Positive Test Scenario with short id enabled",
			id:             "123",
			dbMockFunc:     func() error { return nil },
			shortIdEnabled: true,
			expectedExist:  true,
		},
		{
			name:           "Negative Test Scenario with short id enabled",
			id:             "invalid",
			dbMockFunc:     func() error { return nil },
			shortIdEnabled: true,
			expectedError:  errors.New("invalid id"),
		},
		{
			name:          "Exception Test Scenario with database error",
			id:            "123",
			dbMockFunc:    func() error { return errors.New("database error") },
			expectedError: errors.InternalServer(reason.DatabaseError),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ar.data.DB.Context = func(ctx context.Context) *xorm.Session {
				return &xorm.Session{
					IDFunc: func(id interface{}) *xorm.Session {
						return &xorm.Session{
							GetFunc: func(bean interface{}) (bool, error) {
								if tc.dbMockFunc() != nil {
									return false, tc.dbMockFunc()
								}
								answer := bean.(*entity.Answer)
								answer.ID = "123"
								answer.QuestionID = "456"
								return true, nil
							},
						}
					},
				}
			}

			ctx := context.Background()
			if tc.shortIdEnabled {
				ctx = context.WithValue(ctx, constant.ShortIDFlag, true)
			}

			answer, exist, err := ar.GetAnswer(ctx, tc.id)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedExist, exist)
			if exist {
				assert.Equal(t, tc.id, answer.ID)
				assert.Equal(t, "456", answer.QuestionID)
			}
		})
	}
}
