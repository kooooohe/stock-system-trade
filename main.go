package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"system-trade/system-trade"
)

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

func main() {
	ns, err := targetFiles("./data/6103")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ns)
	vs, err := systemtrade.CandleSticks(ns)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(vs)
	for _, v := range vs {
		fmt.Println(v.Date())
	}
}
