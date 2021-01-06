package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	marginSec = 5
)

func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

// Time2msec require "0:11:22.03" time format. and return Unix time (milli sec)
func Time2msec(timeStamp string) int64 {
	// parse timeStamp 0:38:56.08 -> [0, 38, 56, 8] int
	p := regexp.MustCompile("[:.]").Split(timeStamp, -1)
	var pts [4]int
	for i, v := range p {
		pi, _ := strconv.Atoi(v)
		pts[i] = pi
	}
	//generate Unix time
	t := time.Date(1970, 1, 1, pts[0], pts[1], pts[2], pts[3]*1000000, time.UTC)
	millis := t.UnixNano() / 1000000
	return millis
}

func readVDR(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	lines := make([]string, 0, 30)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		// "0:01:00.27 end" -> "[0:01:00.27]"
		lines = append(lines, strings.Split(sc.Text(), " ")[0])
	}
	if sc.Err() != nil {
		fmt.Fprintf(os.Stderr, "File scan error : %v", err)
	}

	// validate
	if len(lines) == 0 {
		log.Fatal("Can't detect CM.")
		return nil
	}
	return lines
}

func vdrFormating(vdr []string, tsEndTime string) []string {
	marginTime := fmt.Sprintf("0:00:%02d", marginSec)

	slice := vdr
	if Time2msec(slice[0]) == Time2msec(marginTime) {
		slice = slice[1:]
	} else if Time2msec(marginTime) < Time2msec(slice[0]) {
		slice = append([]string{marginTime}, slice...)
	} else {
		for {
			tmp := slice[0]
			slice = slice[1:]
			if Time2msec(marginTime) < Time2msec(tmp) {
				if len(slice)%2 == 0 {
					slice = append([]string{tmp}, slice...)
				} else {
					slice = append([]string{marginTime, tmp}, slice...)
				}
				break
			}
			// validation
			if len(slice) == 0 {
				log.Fatal("slice become null.")
			}
		}
	}

	if len(slice) < 0 {
		slice = append([]string{marginTime}, slice...)
	}

	if abs(Time2msec(tsEndTime)-Time2msec(slice[len(slice)-1])) < 3 {
		slice = slice[:len(slice)-1]
	} else {
		slice = append(slice, tsEndTime)
	}

	// validation
	if len(slice)%2 != 0 {
		log.Fatal("vdr Formating error.")
	}

	return slice
}

func main() {
	recTime := readVDR("./sample.vdr")
	recTime = vdrFormating(recTime, "1:00:01.16")
	for i := 0; i < len(recTime)/2; i++ {
		fmt.Printf("%d : %v to %v\n", i, recTime[i*2], recTime[i*2+1])
	}
}
