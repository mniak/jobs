package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCompositeJob_StartAndWait(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobMock1 := NewMockJob(ctrl)
	jobMock2 := NewMockJob(ctrl)

	compositeJob := CompositeJob{
		Jobs: []Job{
			jobMock1,
			jobMock2,
		},
	}

	stop := make(chan struct{})

	ctxStart := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())
	jobMock1.EXPECT().Start(ctxStart)
	jobMock2.EXPECT().Start(ctxStart)
	compositeJob.Start(ctxStart)

	go func() {
		fakeError1 := errors.New(gofakeit.SentenceSimple())
		jobMock1.EXPECT().
			Wait().
			Return(fakeError1)
		fakeError2 := errors.New(gofakeit.SentenceSimple())
		jobMock2.EXPECT().
			Wait().
			Return(fakeError2)
		compositeJob.Wait()
		close(stop)
	}()

	ctxStop := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())
	fakeError1 := errors.New(gofakeit.SentenceSimple())
	jobMock1.EXPECT().
		Stop(gomock.Any()).
		Return(fakeError1)
	fakeError2 := errors.New(gofakeit.SentenceSimple())
	jobMock2.EXPECT().
		Stop(gomock.Any()).
		Return(fakeError2)
	compositeJob.Stop(ctxStop)
	<-stop
}

func TestCompositeJob_WhenMultipleStartErrors_ShouldGroupAll(t *testing.T) {
	t.Run("When both fail, none should Stop", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fakeError1 := errors.New(gofakeit.SentenceSimple())
		fakeError2 := errors.New(gofakeit.SentenceSimple())

		jobMock1 := NewMockJob(ctrl)
		jobMock1.EXPECT().
			Start(context.TODO()).
			Return(fakeError1)

		jobMock2 := NewMockJob(ctrl)
		jobMock2.EXPECT().
			Start(context.TODO()).
			Return(fakeError2)

		compositeJob := CompositeJob{
			Jobs: []Job{
				jobMock1,
				jobMock2,
			},
		}

		err := compositeJob.Start(context.TODO())
		assert.Error(t, err)
		assert.ErrorIs(t, err, fakeError1)
		assert.ErrorIs(t, err, fakeError2)
	})

	t.Run("When only 1st fails, 2nd should Start and Stop", func(t *testing.T) {
		t.Run("when Stop does not fail", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError1 := errors.New(gofakeit.SentenceSimple())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(context.TODO()).
				Return(fakeError1)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(context.TODO()).
				Return(nil)
			jobMock2.EXPECT().
				Stop(context.TODO()).
				Return(nil)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(context.TODO())
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeError1)
		})
		t.Run("when Stop fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeError1 := errors.New(gofakeit.SentenceSimple())
			fakeErrorStop2 := errors.New(gofakeit.SentenceSimple())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(context.TODO()).
				Return(fakeError1)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(context.TODO()).
				Return(nil)
			jobMock2.EXPECT().
				Stop(context.TODO()).
				Return(fakeErrorStop2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(context.TODO())
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

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(context.TODO()).
				Return(nil)
			jobMock1.EXPECT().
				Stop(context.TODO()).
				Return(nil)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(context.TODO()).
				Return(fakeError2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(context.TODO())
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeError2)
		})

		t.Run("when Stop fails", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fakeErrorStop1 := errors.New(gofakeit.SentenceSimple())
			fakeError2 := errors.New(gofakeit.SentenceSimple())

			jobMock1 := NewMockJob(ctrl)
			jobMock1.EXPECT().
				Start(context.TODO()).
				Return(nil)
			jobMock1.EXPECT().
				Stop(context.TODO()).
				Return(fakeErrorStop1)

			jobMock2 := NewMockJob(ctrl)
			jobMock2.EXPECT().
				Start(context.TODO()).
				Return(fakeError2)

			compositeJob := CompositeJob{
				Jobs: []Job{
					jobMock1,
					jobMock2,
				},
			}

			err := compositeJob.Start(context.TODO())
			assert.Error(t, err)
			assert.ErrorIs(t, err, fakeErrorStop1)
			assert.ErrorIs(t, err, fakeError2)
		})
	})
}
