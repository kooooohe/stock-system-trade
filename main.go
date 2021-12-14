package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CandleStick struct {
	date  time.Time
	start int
	high  int
	low   int
	end   int
}

type DMA struct {
	candleSticks []CandleStick
}

func (d DMA) CurrentAvarage(du int) int {
	return 0
}
func (d DMA) CurrentAvarageWithin(du int, b int) int {
	return 0
}

func targetFiles(tDir string) (tFiles []string, err error) {
	err = filepath.Walk(tDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		// skip to get the target dir self name
		if path == tDir {
			return nil
		}
		tFiles = append(tFiles, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func CandleSticks(paths []string) ([]CandleStick, error) {
	rc := []CandleStick{}
	for _, v := range paths {
		f, err := os.Open(v)
		if err != nil {
			return nil, err
		}
		r := csv.NewReader(f)
		cnt := 2
		r.FieldsPerRecord = -1

		for {
			recode, err := r.Read()
			if cnt > 0 {
				cnt--
				continue
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			// 日付	始値	高値	安値	終値	出来高	終値調整値
			s, _ := strconv.Atoi(recode[1])
			h, _ := strconv.Atoi(recode[2])
			l, _ := strconv.Atoi(recode[3])
			e, _ := strconv.Atoi(recode[4])
			c := CandleStick{
				// TODO date
				start: s,
				high:  h,
				low:   l,
				end:   e,
			}
			fmt.Printf("%T: %v\n", recode, recode)
			rc = append(rc, c)
		}

		f.Close()
	}
	return rc, nil
}

func main() {
	ns, err := targetFiles("./data/6103")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ns)
	vs, err := CandleSticks(ns)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(vs)
	for _, v := range vs {
		fmt.Println(v)
	}
}
