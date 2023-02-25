package twitch

import (
	"errors"
	"github.com/samber/lo"
	"reflect"
	"sync"
)

type chunkedRequestData[Request any, Response any] struct {
	ids                []string
	requestFn          func(Request) (Response, error)
	responseSelectorFn func(Response) interface{}
	paramFn            func(chunk []string) Request
}

func getDataChunked[T, Req, Res any](req *chunkedRequestData[Req, Res]) ([]T, error) {
	results := make([]T, 0, len(req.ids))

	chunkedIds := lo.Chunk(req.ids, 100)

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	errChan := make(chan error, len(chunkedIds))

	for _, chunk := range chunkedIds {
		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()

			data, err := req.requestFn(req.paramFn(chunk))

			if err != nil {
				errChan <- err
				return
			}

			resultValue := reflect.ValueOf(data)

			if reflect.Indirect(resultValue).FieldByName("ErrorMessage").String() != "" {
				errChan <- errors.New(reflect.Indirect(resultValue).FieldByName("ErrorMessage").String())
				return
			}

			selectedField := req.responseSelectorFn(data)

			mu.Lock()
			results = append(
				results,
				selectedField.([]T)...,
			)
			mu.Unlock()
		}(chunk)
	}

	wg.Wait()

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}
