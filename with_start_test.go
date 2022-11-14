package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWithPreStart(t *testing.T) {
	t.Run("When no errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var calls []string

		mockJob := NewMockJob(ctrl)
		job := WithPreStart(mockJob, func(ctx context.Context) error {
			calls = append(calls, "pre start")
			return nil
		})
		mockJob.EXPECT().Start(gomock.Any()).Do(func(_ any) {
			calls = append(calls, "inner start")
		})

		err := job.Start(context.TODO())
		assert.NoError(t, err)
		assert.Len(t, calls, 2)
		assert.Equal(t, "pre start", calls[0])
		assert.Equal(t, "inner start", calls[1])
	})

	t.Run("When pre start fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var calls []string

		fakePreStartError := errors.New(gofakeit.SentenceSimple())

		mockJob := NewMockJob(ctrl)
		job := WithPreStart(mockJob, func(ctx context.Context) error {
			calls = append(calls, "pre start")
			return fakePreStartError
		})

		err := job.Start(context.TODO())
		assert.Equal(t, fakePreStartError, err)
		assert.Len(t, calls, 1)
		assert.Equal(t, "pre start", calls[0])
	})

	t.Run("When inner start fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var calls []string

		fakeInnerStartError := errors.New(gofakeit.SentenceSimple())

		mockJob := NewMockJob(ctrl)
		job := WithPreStart(mockJob, func(ctx context.Context) error {
			calls = append(calls, "pre start")
			return nil
		})
		mockJob.EXPECT().Start(gomock.Any()).DoAndReturn(func(_ any) error {
			calls = append(calls, "inner start")
			return fakeInnerStartError
		})

		err := job.Start(context.TODO())
		assert.Equal(t, fakeInnerStartError, err)
		assert.Len(t, calls, 2)
		assert.Equal(t, "pre start", calls[0])
		assert.Equal(t, "inner start", calls[1])
	})
}
