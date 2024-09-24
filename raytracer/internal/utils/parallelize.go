package utils

import "sync"

func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(f func()) {
			defer waitGroup.Done()
			f()
		}(function)
	}
}
