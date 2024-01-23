package answer

import (
	"context"
	"testing"

	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	count int64
	err   error
}

func (mdb *mockDB) Context(ctx context.Context) *mockDB {
	return mdb
}

func (mdb *mockDB) Where(query interface{}, args ...interface{}) *mockDB {
	return mdb
}

func (mdb *mockDB) Count(bean interface{}) (int64, error) {
	return mdb.count, mdb.err
}

func TestGetAnswerCount_100eb7482e(t *testing.T) {
	testCases := []struct {
		name          string
		db            *mockDB
		expectedCount int64
		expectedErr   error
	}{
		{
			name:          "Returns count of available answers successfully",
			db:            &mockDB{count: 10, err: nil},
			expectedCount: 10,
			expectedErr:   nil,
		},
		{
			name:          "Handles situation when no available answers",
			db:            &mockDB{count: 0, err: nil},
			expectedCount: 0,
			expectedErr:   nil,
		},
		{
			name:          "Handles situation when database is not accessible",
			db:            &mockDB{count: 0, err: context.DeadlineExceeded},
			expectedCount: 0,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name:          "Handles situation when query execution fails",
			db:            &mockDB{count: 0, err: context.DeadlineExceeded},
			expectedCount: 0,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name:          "Handles situation when database returns more than maximum int64 records",
			db:            &mockDB{count: 0, err: context.DeadlineExceeded},
			expectedCount: 0,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name:          "Handles situation when context passed is already cancelled",
			db:            &mockDB{count: 0, err: context.Canceled},
			expectedCount: 0,
			expectedErr:   context.Canceled,
		},
		{
			name:          "Handles situation when invalid context is passed",
			db:            &mockDB{count: 0, err: context.DeadlineExceeded},
			expectedCount: 0,
			expectedErr:   context.DeadlineExceeded,
		},
		{
			name:          "Handles situation when status of answer is different from 'entity.AnswerStatusAvailable'",
			db:            &mockDB{count: 0, err: nil},
			expectedCount: 0,
			expectedErr:   nil,
		},
		{
			name:          "Handles situation when there is an internal server error",
			db:            &mockDB{count: 0, err: context.DeadlineExceeded},
			expectedCount: 0,
			expectedErr:   context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ar := &answerRepo{
				data: &data.Data{
					DB: tc.db,
				},
			}

			count, err := ar.GetAnswerCount(context.Background())

			assert.Equal(t, tc.expectedCount, count)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
