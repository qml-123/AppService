package _func

import (
	"fmt"
	"reflect"
	"time"
)

type RetryEntity struct {
	MaxRetryTimes int
	RetryInterval time.Duration
}

func NewRetryEntity() *RetryEntity {
	return &RetryEntity{
		MaxRetryTimes: 3,
		RetryInterval: 1 * time.Second,
	}
}

func DefaultRetryFunc(fn func(args []interface{}) error, entity *RetryEntity, args ...interface{}) (err error) {
	for i := 0; i < entity.MaxRetryTimes; i++ {
		err = fn(args)
		if err == nil {
			return
		}
		time.Sleep(entity.RetryInterval)
	}
	return
}

func RetryFunc(fn interface{}, entity *RetryEntity, args ...interface{}) (err error) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return fmt.Errorf("expect a function, but got: %v", v.Kind())
	}

	var in []reflect.Value
	for _, param := range args {
		in = append(in, reflect.ValueOf(param))
	}

	for i := 0; i < entity.MaxRetryTimes; i++ {
		result := v.Call(in)

		if len(result) > 0 && result[0].Interface() == nil {
			return nil
		}
		if len(result) > 0 {
			if result[0].Type().ConvertibleTo(reflect.TypeOf(err)) {
				err = result[0].Convert(reflect.TypeOf(err)).Interface().(error)
			} else {
				err = fmt.Errorf("result is not an error")
			}
		}

		time.Sleep(entity.RetryInterval)
	}

	return
}
