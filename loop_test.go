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
	require.NoError(t, looper.Shutdown(ctx1))

	assert.NoError(t, looper.Wait())
	assert.NoError(t, ctx1.Err())
}

func TestStartLoop_WhenStartContextTimeout_ShouldShutdown(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	looper, err := StartLoop(ctx, func(ctx context.Context) {
		time.Sleep(time.Second)
	})
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, context.DeadlineExceeded, looper.Wait())
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestStartLoop_WhenPanics_ShouldNotShutdownLooping(t *testing.T) {
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
		require.NoError(t, looper.Shutdown(ctx))

		assert.NoError(t, looper.Wait())
		assert.GreaterOrEqual(t, count, 10)
	})
}
