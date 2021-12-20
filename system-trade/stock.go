package systemtrade

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	// "strings"
	"time"
)

type CandleStick struct {
	Date  time.Time
	Start int
	High  int
	Low   int
	End   int
}

type CandleSticks []CandleStick

func (c CandleSticks) sort() CandleSticks {
	sort.Slice(c, func(i, j int) bool {
		return c[i].Date.Before(c[j].Date)
	})
	return c
}
func (c CandleSticks) DMA(d int, to int) (r float64) {
	sum := 0
	// size := len(c) -1
	for i := 0; i < d; i++ {
		sum += c[to-i].End
		// println(c[size-i].end)
	}
	return float64(sum) / float64(d)
}

type positionType int

const (
	nothing positionType = 0
	buy     positionType = 1
	sell    positionType = 2
)

type Position struct {
	t     positionType
	price int
	Lc    float64
	Lp    float64 //指値
}

func (po *Position) Buy(p int) {
	po.price = p
	po.t = buy
}
func (po *Position) ShortSell(p int) {
	po.t = sell
	po.price = p
}

func (po *Position) Sell(c CandleStick) (int, float64) {
	// todo あとで
	po.t = nothing

	// LC
	if float64(c.Start) <= float64(po.price)*(1-po.Lc) {
		return c.Start - po.price, -(1.0 - float64(c.Start)/float64(po.price))
	}

	if float64(c.Low) <= float64(po.price)*(1-po.Lc) {
	  return int(math.Ceil(float64(po.price) * po.Lc)), -po.Lc
	}



	return int(math.Ceil(float64(po.price) * po.Lc)), -po.Lc
	// TODO 利確のときも一緒にする？
	// TODO 割合を出す？数値だけだと、具体的な期待値がでない
}

func (po *Position) BuyBack() {
}
func (po Position) IsDoing() bool {
	return po.t != nothing
}
func (po Position) IsBuying() bool {
	return po.t == buy
}
func (po Position) IsSelling() bool {
	return po.t == sell
}
func (po Position) Price() int {
	return po.price
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

func MakeCandleSticks(paths []string) (CandleSticks, error) {
	rc := CandleSticks{}
	for _, v := range paths {
		f, err := os.Open(v)
		if err != nil {
			return nil, err
		}
		r := csv.NewReader(f)
		cnt := 1
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
			// 日付,始値,高値,安値,終値,出来高,調整後終値
			// d := strings.Replace(recode[0], "-", "/", -1)
			layout := "2006/1/2"
			// layout := "2006/01/02"
			dt, _ := time.Parse(layout, recode[0])
			s, _ := strconv.Atoi(recode[1])
			h, _ := strconv.Atoi(recode[2])
			l, _ := strconv.Atoi(recode[3])
			e, _ := strconv.Atoi(recode[4])
			c := CandleStick{
				Date:  dt,
				Start: s,
				High:  h,
				Low:   l,
				End:   e,
			}
			// fmt.Printf("%T: %v\n", recode, recode)
			rc = append(rc, c)
		}

		f.Close()
	}
	return rc.sort(), nil
}
