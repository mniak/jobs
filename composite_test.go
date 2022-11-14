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

	ctxShutdown := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

	fakeErrorShutdown1 := errors.New(gofakeit.SentenceSimple())
	jobMock1.EXPECT().
		Shutdown(ctxShutdown).
		Return(fakeErrorShutdown1)

	fakeErrorShutdown2 := errors.New(gofakeit.SentenceSimple())
	jobMock2.EXPECT().
		Shutdown(ctxShutdown).
		Return(fakeErrorShutdown2)

	errShutdown := compositeJob.Shutdown(ctxShutdown)
	require.Error(t, errShutdown)
	assert.ErrorIs(t, errShutdown, fakeErrorShutdown1)
	assert.ErrorIs(t, errShutdown, fakeErrorShutdown2)
}

func TestCompositeJob_WhenMultipleStartErrors_ShouldGroupAll(t *testing.T) {
	t.Run("When both fail, none should Shutdown", func(t *testing.T) {
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

	t.Run("When only 1st fails, 2nd should Start and Shutdown", func(t *testing.T) {
		t.Run("when Shutdown does not fail", func(t *testing.T) {
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
				Shutdown(ctxStart).
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
		t.Run("when Shutdown fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError1 := errors.New(gofakeit.SentenceSimple())
			fakeErrorShutdown2 := errors.New(gofakeit.SentenceSimple())

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
				Shutdown(ctxStart).
				Return(fakeErrorShutdown2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(ctxStart)
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeError1)
			assert.ErrorIs(t, err, fakeErrorShutdown2)
		})
	})

	t.Run("When 2nd fails, only 1st should Shutdown", func(t *testing.T) {
		t.Run("when Shutdown does not fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError2 := errors.New(gofakeit.SentenceSimple())

			ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(ctxStart).
				Return(nil)
			jobMock1.EXPECT().
				Shutdown(ctxStart).
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

		t.Run("when Shutdown fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeErrorShutdown1 := errors.New(gofakeit.SentenceSimple())
			fakeError2 := errors.New(gofakeit.SentenceSimple())

			ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(ctxStart).
				Return(nil)
			jobMock1.EXPECT().
				Shutdown(ctxStart).
				Return(fakeErrorShutdown1)

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
			assert.ErrorIs(t, err, fakeErrorShutdown1)
			assert.ErrorIs(t, err, fakeError2)
		})
	})
}
