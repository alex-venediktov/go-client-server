package main

import (
	"api"
	"sync"
	"web"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		api.Run(5000)
		wg.Done()
	}()
	go func() {
		web.Run()
		wg.Done()
	}()

	wg.Wait()
}
