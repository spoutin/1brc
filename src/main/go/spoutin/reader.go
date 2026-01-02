package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func (m *Measurement) AddNewData(temp float64) {
	if temp > m.max {
		m.max = temp
	}
	if temp < m.min {
		m.min = temp
	}
	m.cnt++
	m.sum += temp
}

func (m *Measurement) Sum(measurement Measurement) {
	if measurement.max > m.max {
		m.max = measurement.max
	}
	if measurement.min < m.min {
		m.min = measurement.min
	}
	m.sum += measurement.sum
	m.cnt += measurement.cnt
}

func readFile(filename string, startByte, endByte int64) map[string]*Measurement {
	cities := make(map[string]*Measurement)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Find Start of reader
	for startByte > 0 {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, int64(startByte))
		if err != nil {
			panic(err)
		}
		if fmt.Sprintf("%s", b[:n]) == "\n" {
			startByte++
			break // New line character found
		}
		startByte--
	}

	// Find the end of Reader
	fileStat, _ := file.Stat()
	for endByte > startByte && int64(endByte) <= fileStat.Size() {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, int64(endByte))
		if err != nil {
			panic(err)
		}
		if fmt.Sprintf("%s", b[:n]) == "\n" {
			n, _ = file.ReadAt(b, int64(endByte))
			break
		}
		endByte--
	}

	fileScanner := bufio.NewScanner(io.NewSectionReader(file, int64(startByte), int64(endByte-startByte)))
	for fileScanner.Scan() {
		line := fileScanner.Text()
		//fmt.Println(line)
		city, t, _ := strings.Cut(line, ";")
		temp, err := strconv.ParseFloat(t, 64)
		if err != nil {
			panic(err)
		}
		if _, ok := cities[city]; ok {
			cities[city].AddNewData(temp)
		} else {
			cities[city] = &Measurement{
				min: temp,
				max: temp,
				sum: temp,
				cnt: 1,
			}
		}
	}
	return cities
}
