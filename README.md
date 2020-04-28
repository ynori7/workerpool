# Worker Pool
The worker pool library abstracts the setup around creating worker pools, so all
you need to take care of is the actual business logic.

# How it works
The work will be evenly distributed to N workers which process in parallel. Successful
responses are passed back through a success channel and errors through an error channel.
A callback can be specified to provide behavior for successes and errors. The actual logic
of the worker is specified by providing a function which handles the logic and returns either
a result or an error.

# Usage

The worker pool is defined as follow:

```go
func NewWorkerPool(
	workerCount int,
	onSuccess SuccessFunc,
	onError ErrorFunc,
	doWork DoWorkFunc,
)
```
- The workerCount is the number of workers which will process jobs in parallel.
- onSuccess is the callback which is called for each successful result
- onError is the callback which is called for each failed result
- doWork is the function which is called to process each job. This is the part 
which works concurrently.

Here is a short example:

```go
successes := make([]int, 0)

workerPool := NewWorkerPool(
    3, //The number of workers which should work in parallel
    func(result interface{}) { //On success
        r := result.(int)
        successes = append(successes, r) 
    },
    func(err error) { //On error
        log.Println(err.Error())
    },
    func(job interface{}) (result interface{}, err error) { //Do work
        j := job.(int)
        if j > 4 {
            return nil, fmt.Errorf("number too big: %d", j)
        }
        return j, nil
    })

//Do the work
if err := workerPool.Work(ctx, []int{1, 2, 3, 4, 5, 6, 7}); err != nil {
    log.Println(err.Error())
}
```

### Cancelling the jobs

Sometimes you may want to stop processing early if, for example, enough results have 
been found. This can be done by canceling the context passed into the worker pool:

```go
ctx, cancel := context.WithCancel(context.Background())

successes := make([]int, 0)

workerPool := NewWorkerPool(
    3, //The number of workers which should work in parallel
    func(result interface{}) { //On success
        r := result.(int)
        successes = append(successes, r) 
        if len(successes) > 3 {
            cancel() //we have enough results
        }
    },
    func(err error) { //On error
        log.Println(err.Error())
    },
    func(job interface{}) (result interface{}, err error) { //Do work
        j := job.(int)
        if j > 4 {
            return nil, fmt.Errorf("number too big: %d", j)
        }
        return j, nil
    })

//Do the work
if err := workerPool.Work(ctx, []int{1, 2, 3, 4, 5, 6, 7}); err != nil {
    log.Println(err.Error())
}
```
In the above code, once 3 successes have been found, it will signal the worker pool
to 
stop processing futher.
