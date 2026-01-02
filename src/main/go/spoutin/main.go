package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
)

type Measurement struct {
	min float64
	max float64
	cnt int
	sum float64
}

func (m *Measurement) Average() float64 {
	return m.sum / float64(m.cnt)
}

func (m *Measurement) String() string {
	return fmt.Sprintf("Min: %f, Max: %f, Avg: %f", m.min, m.max, m.Average())
}

func Extent(mu *sync.Mutex, original map[string]Measurement, new map[string]*Measurement) {
	mu.Lock()
	defer mu.Unlock()
	for key, value := range new {
		if ci, ok := original[key]; ok {
			ci.Sum(*value)
		} else {
			original[key] = *value
		}
	}
}

func main() {
	measurementFile := flag.String("i", "", "input file")
	threadCount := flag.Int64("t", int64(runtime.NumCPU()), "number of threads")
	flag.Parse()
	if *measurementFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*measurementFile)
	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := fileStat.Size()
	chuckSize := fileSize / *threadCount

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	cities := make(map[string]Measurement, 10_000)
	for i := int64(0); i < *threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := readFile(*measurementFile, chuckSize*i, (chuckSize*i)+chuckSize)
			Extent(mu, cities, m)
		}()
	}
	// Get last chuck
	wg.Add(1)
	func() {
		defer wg.Done()
		m := readFile(*measurementFile, (chuckSize**threadCount)-1, fileStat.Size()-1)
		Extent(mu, cities, m)
	}()

	wg.Wait()

	for city, measurement := range cities {
		fmt.Printf("%s - %s\n", city, measurement.String())
	}

}
