package answer

import (
	"context"
	"errors"
	"testing"

	"github.com/apache/incubator-answer/internal/base/data"
	"github.com/apache/incubator-answer/internal/entity"
	"github.com/apache/incubator-answer/internal/service/activity_common"
	"github.com/apache/incubator-answer/internal/service/rank"
	"github.com/apache/incubator-answer/internal/service/unique"
	"github.com/stretchr/testify/assert"
)

type mockUniqueIDRepo struct {
	unique.UniqueIDRepo
}

type mockUserRankRepo struct {
	rank.UserRankRepo
}

type mockActivityRepo struct {
	activity_common.ActivityRepo
}

func (m *mockUniqueIDRepo) GenUniqueIDStr(ctx context.Context, key string) (uniqueID string, err error) {
	return "", nil
}

func (m *mockUserRankRepo) GetMaxDailyRank(ctx context.Context) (maxDailyRank int, err error) {
	return 0, nil
}

func (m *mockActivityRepo) GetActivityTypeByObjID(ctx context.Context, objectId string, action string) (activityType, rank int, hasRank int, err error) {
	return 0, 0, 0, nil
}

func TestAddAnswer_513bc5b791(t *testing.T) {
	testCases := []struct {
		desc        string
		answer      *entity.Answer
		uniqueIDErr error
		wantErr     error
	}{
		{
			desc:   "success",
			answer: &entity.Answer{QuestionID: "1234567890"},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockUniqueIDRepo := &mockUniqueIDRepo{}
			mockUserRankRepo := &mockUserRankRepo{}
			mockActivityRepo := &mockActivityRepo{}

			ar := &answerRepo{
				data:         &data.Data{},
				uniqueIDRepo: mockUniqueIDRepo,
				userRankRepo: mockUserRankRepo,
				activityRepo: mockActivityRepo,
			}

			err := ar.AddAnswer(context.Background(), tC.answer)

			if tC.wantErr != nil {
				assert.EqualError(t, err, tC.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tC.answer.ID)
				assert.NotEmpty(t, tC.answer.QuestionID)
			}
		})
	}
}
