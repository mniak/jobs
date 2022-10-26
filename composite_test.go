package jobs

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
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
		ctxWait := context.WithValue(context.TODO(), gofakeit.Word(), gofakeit.Word())
		fakeError1 := errors.New(gofakeit.SentenceSimple())
		jobMock1.EXPECT().
			Wait(gomock.Any()).
			Return(fakeError1)
		fakeError2 := errors.New(gofakeit.SentenceSimple())
		jobMock2.EXPECT().
			Wait(gomock.Any()).
			Return(fakeError2)
		compositeJob.Wait(ctxWait)
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
