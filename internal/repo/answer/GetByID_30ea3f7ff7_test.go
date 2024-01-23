package answer

import (
	"context"
	"testing"

	"github.com/apache/incubator-answer/internal/base/constant"
	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/pkg/uid"
	"github.com/segmentfault/pacman/errors"
)

func TestGetByID_30ea3f7ff7(t *testing.T) {
	mockAnswer := &entity.Answer{
		ID:         "1234567890",
		QuestionID: "0987654321",
	}

	mockData := &data.Data{
		DB: new(MockDB),
	}

	ar := &answerRepo{
		data: mockData,
	}

	tests := []struct {
		name      string
		answerID  string
		want      *entity.Answer
		wantFound bool
		wantErr   error
		shortID   bool
	}{
		{
			name:      "Valid Answer ID",
			answerID:  "1234567890",
			want:      mockAnswer,
			wantFound: true,
			wantErr:   nil,
			shortID:   false,
		},
		{
			name:      "Invalid Answer ID",
			answerID:  "invalid",
			want:      nil,
			wantFound: false,
			wantErr:   nil,
			shortID:   false,
		},
		{
			name:      "Empty Answer ID",
			answerID:  "",
			want:      nil,
			wantFound: false,
			wantErr:   nil,
			shortID:   false,
		},
		{
			name:      "Database Error",
			answerID:  "1234567890",
			want:      nil,
			wantFound: false,
			wantErr:   errors.InternalServer("DatabaseError"),
			shortID:   false,
		},
		{
			name:     "Shortened Answer ID",
			answerID: "1234567890",
			want: &entity.Answer{
				ID:         uid.EnShortID("1234567890"),
				QuestionID: uid.EnShortID("0987654321"),
			},
			wantFound: true,
			wantErr:   nil,
			shortID:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.shortID {
				ctx = context.WithValue(ctx, constant.ShortIDFlag, true)
			}

			got, gotFound, gotErr := ar.GetByID(ctx, tt.answerID)
			if got != tt.want {
				t.Errorf("GetByID() got = %v, want = %v", got, tt.want)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetByID() gotFound = %v, want = %v", gotFound, tt.wantFound)
			}
			if gotErr != tt.wantErr {
				t.Errorf("GetByID() gotErr = %v, want = %v", gotErr, tt.wantErr)
			}
		})
	}
}

type MockDB struct {
}

func (db *MockDB) Context(ctx context.Context) *MockDB {
	return db
}

func (db *MockDB) ID(id interface{}) *MockDB {
	return db
}

func (db *MockDB) Get(bean interface{}) (bool, error) {
	if id, ok := bean.(*entity.Answer); ok {
		if id.ID == "1234567890" {
			return true, nil
		}
		if id.ID == "" {
			return false, nil
		}
	}
	return false, errors.InternalServer("DatabaseError")
}
