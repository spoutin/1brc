package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Cities sync.Map

var (
	readChan = make(chan string, 100)
)

type Measurement struct {
	mu  sync.Mutex
	min float32
	max float32
	avg float32
	cnt int
}

func (m *Measurement) String() string {
	return fmt.Sprintf("%f,%f,%f", m.min, m.max, m.avg)
}

func (m *Measurement) AddNewData(temp float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if temp > m.max {
		m.max = temp
	}
	if temp < m.min {
		m.min = temp
	}
	m.cnt++
	// Calculate avg without storing all values
	// Avg = Old Avg + (New Value - Old Avg) / New Count
	m.avg = (m.avg + (temp - m.avg)) / float32(m.cnt)
}

func startReader(measurementFile string, startLine, endLine int, readChan chan string) {
	file, err := os.Open(measurementFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	currentLine := startLine
	slog.Info(fmt.Sprintf("Starting reader for %d", currentLine))
	for scanner.Scan() {
		if currentLine >= startLine && currentLine <= endLine {
			readChan <- scanner.Text()
		} else if currentLine > endLine {
			// Stop reading
			break
		}
		currentLine++
	}
}

func worker(cities *sync.Map, readChan chan string, closeChan chan bool) {
	for {
		select {
		case <-closeChan:
			return
		case data := <-readChan:
			measurement := strings.Split(data, ";")
			f, err := strconv.ParseFloat(measurement[1], 32)
			if err != nil {
				panic(err)
			}
			if d, loaded := cities.LoadOrStore(measurement[0], &Measurement{
				min: float32(f),
				max: float32(f),
				avg: float32(f),
				cnt: 1,
			}); loaded {
				d.(*Measurement).AddNewData(float32(f))
			}
		}
	}
}

func main() {
	measurementFile := flag.String("i", "", "input file")
	flag.Parse()
	if *measurementFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	closeChannel := make(chan bool)
	var cities sync.Map
	wg := new(sync.WaitGroup)
	wgReader := new(sync.WaitGroup)
	defer close(readChan)
	for i := 0; i < 10; i++ {
		wgReader.Add(1)
		go func() {
			defer wgReader.Done()
			startReader(*measurementFile, i*1_000_000+i, (i*1_000_000)+1_0000_000, readChan)
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(&cities, readChan, closeChannel)
		}()
	}
	wgReader.Wait()
	close(closeChannel) // Add Readers are complete
	wg.Wait()
	cities.Range(func(key, value interface{}) bool {
		fmt.Printf("city %s, min: %f, max: %f, avg: %f, count: %d\n",
			key, value.(*Measurement).min,
			value.(*Measurement).min,
			value.(*Measurement).max,
			value.(*Measurement).cnt,
		)
		return true
	})
}
