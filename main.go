package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	systemtrade "system-trade/system-trade"
	"time"
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

func trade(sDate time.Time, cs systemtrade.CandleSticks) {
	skipc := 10
	p := 0

	po := systemtrade.Position{Lc: 0.03, Lp: 0.07}
	for i, v := range cs {
		// for DMA
		if skipc > 0 || v.Date.Before(sDate) {
			skipc--
			continue
		}
		d := cs.DMA(10, i)
		fmt.Printf("%v: ", v.Date)
		fmt.Println(d)
		// fmt.Println(v.Date())
		//TODO 10DMAが上向きで、株価のしたひげでも一回でもDMA以下にあって、次の日が高値を超えたら 3% 7%

		wasDMAUp := cs.DMA(10, i-1) > cs.DMA(10, i-2)
		wasStockUnderDMA := float64(v.Low) < cs.DMA(10, i-1)
		wasStockOverDMA := float64(v.High) > cs.DMA(10, i-1)
		yesterday := cs[i-1]

		if po.IsBuying() {
			if float64(v.Low) <= float64(po.Price()) * (1 -po.Lc) {
			}
			continue
		}
		if po.IsSelling() {

			continue
		}

		if wasDMAUp && wasStockUnderDMA {
			if v.High > yesterday.High {
				p = yesterday.High + 1
				if v.Start > yesterday.High {
					p = v.Start
				}
				po.Buy(p)
				fmt.Printf("Buy: %v: %v", v.Date, p)
			}
		}

		if !wasDMAUp && wasStockOverDMA {
			if v.Low < yesterday.Low {
				p = yesterday.Low - 1
				if v.Start < yesterday.Low {
					p = v.Start
				}
				po.ShortSell(p)
				fmt.Printf("ShortSell: %v: %v", v.Date, p)
			}
		}

	}
	// vs.DMA(10)
}

func main() {
	ns, err := targetFiles("./data/6103")
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(ns)
	vs, err := systemtrade.MakeCandleSticks(ns)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(vs)
	// cnt := 0
	layout := "2006/01/02"
	sDate, _ := time.Parse(layout, "2015/01/01")
	trade(sDate, vs)
}
