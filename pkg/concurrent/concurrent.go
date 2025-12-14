package concurrent

import "context"

type ProcessResult[T any] struct {
	Result T
	Error  error
	Index  int
}

func ProcessItems[T any, R any](
	ctx context.Context,
	items []T,
	maxConcurrent int,
	processFunc func(context.Context, T) (R, error),
) ([]R, []error) {
	if len(items) == 0 {
		return nil, nil
	}

	resultsChan := make(chan ProcessResult[R], len(items))
	semaphore := make(chan struct{}, maxConcurrent)

	for i, item := range items {
		i, item := i, item
		semaphore <- struct{}{}

		go func() {
			defer func() { <-semaphore }()

			result, err := processFunc(ctx, item)
			resultsChan <- ProcessResult[R]{
				Result: result,
				Error:  err,
				Index:  i,
			}
		}()
	}

	results := make([]R, len(items))
	errors := make([]error, 0)

	for i := 0; i < len(items); i++ {
		result := <-resultsChan
		if result.Error != nil {
			errors = append(errors, result.Error)
		} else {
			results[result.Index] = result.Result
		}
	}

	filteredErrors := make([]error, 0)
	for _, err := range errors {
		if err != nil {
			filteredErrors = append(filteredErrors, err)
		}
	}

	return results, filteredErrors
}
