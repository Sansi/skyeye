package main

import (
	"bytes"
	"fmt"
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
	numBenchmark := 500

	testIds := make([]string, numBenchmark)
	for i := 0; i < len(testIds); i++ {
		testIds[i] = "test" + randomString(11)
	}

	client.NewDevices(testIds)
	defer client.DelDevices(testIds)

	for i := 0; i < numBenchmark; i++ {
		fmt.Println("CreateClient: ", testIds[i])
		go client.CreateClient(cid, testIds[i])
		time.Sleep(300 * time.Millisecond)
	}
	for i := 0; i < numBenchmark; i++ {
		<-cid
	}

}
