package workerpool

import (
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
	workerPool := NewWorkerPool(3,
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
	err := workerPool.Work([]int{1, 2, 3, 4, 5, 6, 7})

	//then
	require.NoError(t, err, "there should not have been an error")
	assert.Equal(t, 4, successCount)
	assert.Equal(t, 3, errorCount)
}

func Test_DoWork_InputNotSlice(t *testing.T) {
	//given
	workerPool := NewWorkerPool(3,
		func(result interface{}) {
		},
		func(err error) {
		},
		func(job interface{}) (result interface{}, err error) {
			return nil, nil
		})

	//when
	err := workerPool.Work(5)

	//then
	assert.Error(t, err, "there should have been an error")
}
