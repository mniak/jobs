package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartLoop_HappyScenario(t *testing.T) {
	done := make(chan struct{})
	var count int

	ctx1, _ := context.WithTimeout(context.Background(), 2*time.Second)
	looper, err := StartLoop(ctx1, func(ctx context.Context) {
		if count == 5 {
			close(done)
		}
		count++
	})
	require.NoError(t, err)

	<-done
	require.NoError(t, looper.Stop(ctx1))

	ctx2, _ := context.WithTimeout(context.Background(), 2*time.Second)
	looper.Wait(ctx2)

	assert.NoError(t, ctx1.Err())
	assert.NoError(t, ctx2.Err())
}

func TestStartLoop_WhenStartContextTimeout_ShouldStop(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	looper, err := StartLoop(ctx, func(ctx context.Context) {
		time.Sleep(time.Second)
	})
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, context.DeadlineExceeded, looper.Wait(context.Background()))
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestStartLoop_WhenWaitTimeout_ShouldStop(t *testing.T) {
	looper, err := StartLoop(context.Background(), func(ctx context.Context) {
		time.Sleep(60 * time.Second)
	})
	require.NoError(t, err)

	ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)

	err = looper.Wait(ctx)
	assert.EqualError(t, err, "wait: context deadline exceeded")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestStartLoop_WhenPanics_ShouldNotStopLooping(t *testing.T) {
	done := make(chan struct{})
	var count int

	assert.NotPanics(t, func() {
		ctx := context.Background()
		looper, err := StartLoop(ctx, func(ctx context.Context) {
			switch count {
			case 5:
				count++
				panic("this is the panic 5")
			case 10:
				close(done)
			default:
				count++
			}
		})
		require.NoError(t, err)

		<-done
		require.NoError(t, looper.Stop(ctx))

		looper.Wait(ctx)
		assert.NoError(t, ctx.Err())
		assert.GreaterOrEqual(t, count, 10)
	})
}
