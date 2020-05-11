package workerpool

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DoWork(t *testing.T) {
	//given
	expectedErr := fmt.Errorf("something went wrong")

	errorCount := 0
	successCount := 0
	workerPool := NewWorkerPool(
		func(result interface{}) {
			r := result.(int)
			assert.Equal(t, true, r <= 4)
			successCount++
		},
		func(err error) {
			assert.Error(t, expectedErr, err)
			errorCount++
		},
		func(job interface{}) (result interface{}, err error) {
			j := job.(int)
			if j > 4 {
				return nil, expectedErr
			}
			return j, nil
		})

	//when
	err := workerPool.Work(context.Background(), 3, []int{1, 2, 3, 4, 5, 6, 7})

	//then
	require.NoError(t, err, "there should not have been an error")
	assert.Equal(t, 4, successCount)
	assert.Equal(t, 3, errorCount)
}

func Test_DoWork_Error(t *testing.T) {
	//given
	testcases := map[string]struct {
		workerCount int
		jobs        interface{}
		expectedErr error
	}{
		"input not slice": {
			workerCount: 1,
			jobs:        5,
			expectedErr: fmt.Errorf("input is not a slice"),
		},
		"zero workers": {
			workerCount: 0,
			jobs:        []int{5},
			expectedErr: fmt.Errorf("there must be at least one worker"),
		},
		"negative workers": {
			workerCount: -1,
			jobs:        []int{5},
			expectedErr: fmt.Errorf("there must be at least one worker"),
		},
	}

	workerPool := NewWorkerPool(
		func(result interface{}) {
		},
		func(err error) {
		},
		func(job interface{}) (result interface{}, err error) {
			return nil, nil
		})

	for testcase, testdata := range testcases {
		//when
		err := workerPool.Work(context.Background(), testdata.workerCount, testdata.jobs)

		//then
		require.Error(t, err, "there should have been an error", testcase)
		assert.EqualError(t, err, testdata.expectedErr.Error())
	}
}

func Test_DoWork_WithCancellation(t *testing.T) {
	//given
	expectedErr := fmt.Errorf("something went wrong")

	errorCount := 0
	successCount := 0
	ctx, cancel := context.WithCancel(context.Background())

	workerPool := NewWorkerPool(
		func(result interface{}) {
			if successCount >= 1 {
				cancel()
			}
			r := result.(int)
			assert.Equal(t, true, r <= 4)
			successCount++
		},
		func(err error) {
			assert.Error(t, expectedErr, err)
			errorCount++
		},
		func(job interface{}) (result interface{}, err error) {
			j := job.(int)
			if j > 4 {
				return nil, expectedErr
			}
			return j, nil
		})

	//when
	err := workerPool.Work(ctx, 1, []int{1, 2, 3, 4, 5, 6, 7, 8})

	//then
	require.NoError(t, err, "there should not have been an error")
	assert.True(t, successCount < 4)
	assert.Equal(t, 0, errorCount, "not all the items should have been processed")
}
