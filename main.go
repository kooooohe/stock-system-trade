package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type CandleStick struct {
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
		return  nil, err
	}
	return 
}

func main() {
	// tDir := "./data/6103"
	ns, err:= targetFiles("./data/6103")
	fmt.Println(ns)
	return

	f, err := os.Open("./data/6103/6103_2015.csv")
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
		fmt.Printf("%T: %v\n", recode, recode)
	}
	// println("test")
}
