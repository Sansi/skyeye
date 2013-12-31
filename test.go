package main

import (
	"bytes"
	// "fmt"
	"github.com/edwardtoday/skyeye/client"
	// "github.com/edwardtoday/skyeye/utils"
	"math/rand"
	"time"
)

func randomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(randInt(48, 57)) != temp {
			temp = string(randInt(48, 57))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func main() {
	cid := make(chan int)
	numBenchmark := 200
	numRun := 10

	testIds := make([]string, numBenchmark*numRun)
	for i := 0; i < len(testIds); i++ {
		testIds[i] = "test" + randomString(11)
	}

	client.NewDevices(testIds)
	defer client.DelDevices(testIds)

	for j := 0; j < numRun; j++ {
		for i := 0; i < numBenchmark; i++ {
			go client.CreateClient(cid, testIds[i+j*numBenchmark])
			time.Sleep(10 * time.Millisecond)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	for i := 0; i < numBenchmark*numRun; i++ {
		<-cid
	}

}
