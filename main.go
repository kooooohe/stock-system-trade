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

func main() {
	subDirToSkip := ""
	err := filepath.Walk("./data/6103", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
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
