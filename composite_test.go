package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompositeJob_HappyScenario(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobMock1 := NewMockJob(ctrl)
	jobMock2 := NewMockJob(ctrl)

	ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

	jobMock1.EXPECT().Start(ctxStart)
	jobMock2.EXPECT().Start(ctxStart)

	compositeJob := CompositeJob{
		Jobs: []Job{
			jobMock1,
			jobMock2,
		},
	}

	err := compositeJob.Start(ctxStart)
	require.NoError(t, err)

	fakeErrorWait1 := errors.New(gofakeit.SentenceSimple())
	jobMock1.EXPECT().
		Wait().
		Return(fakeErrorWait1)

	fakeErrorWait2 := errors.New(gofakeit.SentenceSimple())
	jobMock2.EXPECT().
		Wait().
		Return(fakeErrorWait2)

	errWait := compositeJob.Wait()
	require.Error(t, errWait)
	assert.ErrorIs(t, errWait, fakeErrorWait1)
	assert.ErrorIs(t, errWait, fakeErrorWait2)

	ctxStop := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

	fakeErrorStop1 := errors.New(gofakeit.SentenceSimple())
	jobMock1.EXPECT().
		Stop(ctxStop).
		Return(fakeErrorStop1)

	fakeErrorStop2 := errors.New(gofakeit.SentenceSimple())
	jobMock2.EXPECT().
		Stop(ctxStop).
		Return(fakeErrorStop2)

	errStop := compositeJob.Stop(ctxStop)
	require.Error(t, errStop)
	assert.ErrorIs(t, errStop, fakeErrorStop1)
	assert.ErrorIs(t, errStop, fakeErrorStop2)
}

func TestCompositeJob_WhenMultipleStartErrors_ShouldGroupAll(t *testing.T) {
	t.Run("When both fail, none should Stop", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError1 := errors.New(gofakeit.SentenceSimple())
		fakeError2 := errors.New(gofakeit.SentenceSimple())

		ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

		jobMock1 := NewMockJob(ctrl)
		jobMock1.EXPECT().
			Start(ctxStart).
			Return(fakeError1)

		jobMock2 := NewMockJob(ctrl)
		jobMock2.EXPECT().
			Start(ctxStart).
			Return(fakeError2)

		compositeJob := CompositeJob{
			Jobs: []Job{
				jobMock1,
				jobMock2,
			},
		}

		err := compositeJob.Start(ctxStart)
		assert.Error(t, err)
		assert.ErrorIs(t, err, fakeError1)
		assert.ErrorIs(t, err, fakeError2)
	})

	t.Run("When only 1st fails, 2nd should Start and Stop", func(t *testing.T) {
		t.Run("when Stop does not fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError1 := errors.New(gofakeit.SentenceSimple())

			ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(ctxStart).
				Return(fakeError1)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(ctxStart).
				Return(nil)
			jobMock2.EXPECT().
				Stop(ctxStart).
				Return(nil)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(ctxStart)
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeError1)
		})
		t.Run("when Stop fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError1 := errors.New(gofakeit.SentenceSimple())
			fakeErrorStop2 := errors.New(gofakeit.SentenceSimple())

			ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(ctxStart).
				Return(fakeError1)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(ctxStart).
				Return(nil)
			jobMock2.EXPECT().
				Stop(ctxStart).
				Return(fakeErrorStop2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(ctxStart)
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeError1)
			assert.ErrorIs(t, err, fakeErrorStop2)
		})
	})

	t.Run("When 2nd fails, only 1st should Stop", func(t *testing.T) {
		t.Run("when Stop does not fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError2 := errors.New(gofakeit.SentenceSimple())

			ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(ctxStart).
				Return(nil)
			jobMock1.EXPECT().
				Stop(ctxStart).
				Return(nil)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(ctxStart).
				Return(fakeError2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(ctxStart)
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeError2)
		})

		t.Run("when Stop fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeErrorStop1 := errors.New(gofakeit.SentenceSimple())
			fakeError2 := errors.New(gofakeit.SentenceSimple())

			ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(ctxStart).
				Return(nil)
			jobMock1.EXPECT().
				Stop(ctxStart).
				Return(fakeErrorStop1)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(ctxStart).
				Return(fakeError2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(ctxStart)
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeErrorStop1)
			assert.ErrorIs(t, err, fakeError2)
		})
	})
}
