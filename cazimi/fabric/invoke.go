package fabric

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
	secAction "github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/action"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/chaincode/invoketask"
	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/executor"
	"github.com/snowdiceX/dawns.world/cazimi/fabric/config"
	"github.com/snowdiceX/dawns.world/cazimi/log"
)

// ChaincodeInvoke invoke chaincode
func ChaincodeInvoke(chainID, chaincodeID string, argsArray []Args) (result string, err error) {
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return
	}
	action, err := newChaincodeInvokeAction()
	if err != nil {
		log.Errorf("Error while initializing invokeAction: %v", err)
		return
	}

	defer action.Terminate()

	err = action.invoke(chainID, chaincodeID, argsArray)
	if err != nil {
		log.Errorf("Error while calling action.invoke(): %v", err)
	}
	return
}

type chaincodeInvokeAction struct {
	Action
	numInvoked uint32
	done       chan bool
}

func newChaincodeInvokeAction() (*chaincodeInvokeAction, error) {
	action := &chaincodeInvokeAction{done: make(chan bool)}
	err := action.Initialize()
	return action, err
}

func (a *chaincodeInvokeAction) invoke(chainID, chaincodeID string, argsArray []Args) error {
	channelClient, err := a.ChannelClient()
	if err != nil {
		return errors.Errorf("Error getting channel client: %v", err)
	}

	executor := executor.NewConcurrent("Invoke Chaincode", config.Config().Concurrency)
	executor.Start()
	defer executor.Stop(true)

	success := 0
	var errs []error
	var successDurations []time.Duration
	var failDurations []time.Duration

	var targets []fab.Peer
	if len(config.Config().PeerURL) > 0 || len(config.Config().OrgIDs) > 0 {
		targets = a.Peers()
	}

	var wg sync.WaitGroup
	var mutex sync.RWMutex
	var tasks []*invoketask.Task
	var taskID int
	argsStructArray := transform(argsArray)
	for i := 0; i < config.Config().Iterations; i++ {
		for _, args := range argsStructArray {
			taskID++
			var startTime time.Time
			task := invoketask.New(
				strconv.Itoa(taskID), channelClient, targets,
				chaincodeID,
				&args, executor,
				retry.Opts{
					Attempts:       config.Config().MaxAttempts,
					InitialBackoff: config.Config().InitialBackoff,
					MaxBackoff:     config.Config().MaxBackoff,
					BackoffFactor:  config.Config().BackoffFactor,
					RetryableCodes: retry.ChannelClientRetryableCodes,
				},
				config.Config().Verbose || config.Config().Iterations == 1,
				config.Config().PrintPayloadOnly, a.Printer(),

				func() {
					startTime = time.Now()
				},
				func(err error) {
					duration := time.Since(startTime)
					defer wg.Done()
					mutex.Lock()
					defer mutex.Unlock()
					if err != nil {
						errs = append(errs, err)
						failDurations = append(failDurations, duration)
					} else {
						success++
						successDurations = append(successDurations, duration)
					}
				})
			tasks = append(tasks, task)
		}
	}

	numInvocations := len(tasks)

	wg.Add(numInvocations)

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				mutex.RLock()
				if len(errs) > 0 {
					fmt.Printf("*** %d failed invocation(s) out of %d\n", len(errs), numInvocations)
				}
				fmt.Printf("*** %d successfull invocation(s) out of %d\n", success, numInvocations)
				mutex.RUnlock()
			case <-done:
				return
			}
		}
	}()

	startTime := time.Now()

	for _, task := range tasks {
		if err := executor.Submit(task); err != nil {
			return errors.Errorf("error submitting task: %s", err)
		}
	}

	// Wait for all tasks to complete
	wg.Wait()
	done <- true

	duration := time.Now().Sub(startTime)

	var allErrs []error
	var attempts int
	for _, task := range tasks {
		attempts = attempts + task.Attempts()
		if task.LastError() != nil {
			allErrs = append(allErrs, task.LastError())
		}
	}

	if len(errs) > 0 {
		fmt.Printf("\n*** %d errors invoking chaincode:\n", len(errs))
		for _, err := range errs {
			fmt.Printf("%s\n", err)
		}
	} else if len(allErrs) > 0 {
		fmt.Printf("\n*** %d transient errors invoking chaincode:\n", len(allErrs))
		for _, err := range allErrs {
			fmt.Printf("%s\n", err)
		}
	}

	if numInvocations > 1 {
		fmt.Printf("\n")
		fmt.Printf("*** ---------- Summary: ----------\n")
		fmt.Printf("***   - Invocations:     %d\n", numInvocations)
		fmt.Printf("***   - Concurrency:     %d\n", config.Config().Concurrency)
		fmt.Printf("***   - Successfull:     %d\n", success)
		fmt.Printf("***   - Total attempts:  %d\n", attempts)
		fmt.Printf("***   - Duration:        %2.2fs\n", duration.Seconds())
		fmt.Printf("***   - Rate:            %2.2f/s\n", float64(numInvocations)/duration.Seconds())
		fmt.Printf("***   - Average:         %2.2fs\n", average(append(successDurations, failDurations...)))
		fmt.Printf("***   - Average Success: %2.2fs\n", average(successDurations))
		fmt.Printf("***   - Average Fail:    %2.2fs\n", average(failDurations))
		fmt.Printf("***   - Min Success:     %2.2fs\n", min(successDurations))
		fmt.Printf("***   - Max Success:     %2.2fs\n", max(successDurations))
		fmt.Printf("*** ------------------------------\n")
	}

	return nil
}

func average(durations []time.Duration) float64 {
	if len(durations) == 0 {
		return 0
	}

	var total float64
	for _, duration := range durations {
		total += duration.Seconds()
	}
	return total / float64(len(durations))
}

func min(durations []time.Duration) float64 {
	min, _ := minMax(durations)
	return min
}

func max(durations []time.Duration) float64 {
	_, max := minMax(durations)
	return max
}

func minMax(durations []time.Duration) (min float64, max float64) {
	for _, duration := range durations {
		if min == 0 || min > duration.Seconds() {
			min = duration.Seconds()
		}
		if max == 0 || max < duration.Seconds() {
			max = duration.Seconds()
		}
	}
	return
}

func transform(argsArray []Args) []secAction.ArgStruct {
	var array []secAction.ArgStruct
	for _, a := range argsArray {
		args := secAction.ArgStruct{
			Func: a.Func, Args: a.Args}
		array = append(array, args)
	}
	return array
}
