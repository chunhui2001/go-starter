package actions

import (
	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"
	"time"
)

func ParallelLoop0Handler(c *gin.Context) {

	N := 10
	data := make([]float64, 0, N)

	for i := 0; i < N; i++ {
		data = append(data, float64(i))
	}

	for i := range data {
		data[i] += 0.08002
		time.Sleep(1 * time.Second)
	}

	c.JSON(200, (&R{Data: data}).Success())

}

func ParallelLoop1Handler(c *gin.Context) {

	N := 10000
	data := make([]float64, 0, N)

	for i := 0; i < N; i++ {
		data = append(data, float64(i))
	}

	res := make([]float64, N)
	sem := make(chan utils.Empty, N) // semaphore pattern

	for i, xi := range data {
		go func(i int, xi float64) {
			res[i] = xi + 0.05
			time.Sleep(1 * time.Second)
			sem <- utils.Empty{}
		}(i, xi)
	}

	// wait for goroutines to finish
	for i := 0; i < N; i = i + 1 {
		<-sem
	}

	c.JSON(200, (&R{Data: res}).Success())

}

func ParallelLoop2Handler(c *gin.Context) {

	N := 10
	data := make([]float64, 0, N)

	for i := 0; i < N; i++ {
		data = append(data, float64(i))
	}

	sem := make(utils.Semaphore, len(data))

	for i := range data {
		go func(i int) {
			data[i] += 0.061
			time.Sleep(1 * time.Second)
			sem.Signal()
		}(i)
	}

	sem.Wait(len(data))

	c.JSON(200, (&R{Data: data}).Success())

}

// func ParallelLoop3Handler(c *gin.Context) {

// 	N := 10
// 	data := make([]float64, 0, N)

// 	for i := 0; i < N; i++ {
// 		data = append(data, float64(i))
// 	}

// 	for _, xi := range data {
// 		xch := make(chan float64)
// 		go func() {
// 			xi := <-xch
// 			out <- func(xi float64) float {
// 				return xi + 0.8999111044
// 			}()
// 		}()
// 		xch <- xi
// 	}

// 	c.JSON(200, (&R{Data: data}).Success())

// }
