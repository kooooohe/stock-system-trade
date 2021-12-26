package systemtrade

import (
	"encoding/csv"
	"fmt"
	_ "fmt"
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

// ====== //

type recode struct {
	start
	end
	t          positionType
	defference float64
}

func (r *recode) SetEnd(c CandleStick, p int, d float64) {
	r.end = end{p: p, c: c}
	r.defference = d
}

type start struct {
	c CandleStick
	p int
}

type end struct {
	c CandleStick
	p int
}

type Score struct {
	win       int
	lose      int
	sum       float64
	buy       int
	shortSell int
	recodes   []recode
}

func (s *Score) SetStartRecode(c CandleStick, p int, t positionType) {
	r := recode{start: start{c: c, p: p}, t: t}
	s.recodes = append(s.recodes, r)

	if t == buy {
		s.Buy()
	} else {
		s.ShortSell()
	}
}

func (s *Score) SetEndRcode(c CandleStick, d float64) {
	s.Sum(d)
	//TODO endにpriceはいる？
	s.recodes[len(s.recodes)-1].SetEnd(c,0, d)
}

func (s *Score) Win() {
	s.win++
}

func (s *Score) Lose() {
	s.lose++
}
func (s *Score) Buy() {
	s.buy++
}
func (s *Score) ShortSell() {
	s.shortSell++
}
func (s *Score) Sum(r float64) {
	s.sum += r
}
func (s Score) Out() {
	// fmt.Printf("%#v", s)
	fmt.Printf("sum: %v\n", s.sum)
	fmt.Printf("win %v\n", s.win)
	fmt.Printf("lose %v\n", s.lose)
	fmt.Printf("buy %v\n", s.buy)
	fmt.Printf("shortSell %v\n", s.shortSell)
}

var Result Score

// ====== //

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

func (po *Position) Buy(p int, c CandleStick) {
	po.t = buy

	Result.SetStartRecode(c, p, buy)
}

func (po *Position) ShortSell(p int, c CandleStick) {
	po.t = sell
	po.price = p

	Result.SetStartRecode(c, p, sell)
}

func (po *Position) Sell(c CandleStick) (int, float64, bool) {
	a, b, ok := po.sell(c)
	if ok {
		po.t = nothing
		//println("kohe")
		//println(b)
		Result.SetEndRcode(c,b)
		if b > 0 {
			Result.Win()
		} else {

			Result.Lose()
		}
	}
	return a, b, ok
}

func (po *Position) BuyBack(c CandleStick) (int, float64, bool) {
	a, b, ok := po.buyBack(c)
	if ok {
		//println("kohe")
		// println(b)
		Result.SetEndRcode(c,b)
		po.t = nothing
		if b > 0 {
			Result.Win()
		} else {

			Result.Lose()
		}
	}
	return a, b, ok
}
func (po *Position) sell(c CandleStick) (int, float64, bool) {

	// LC
	if float64(c.Start) <= po.lossCutPrice() {
		r := -(1.0 - float64(c.Start)/float64(po.price))
		return c.Start - po.price, r, true
	}

	// LC
	if float64(c.Low) <= po.lossCutPrice() {
		r := -po.Lc
		return int(math.Ceil(float64(po.price) * po.Lc)), r, true
	}

	// PROFIT
	if float64(c.Start) >= po.profitPrice() {
		r := (float64(c.Start)/float64(po.price) - 1)
		return c.Start - po.price, r, true
	}

	// PROFIT
	if float64(c.High) >= po.profitPrice() {
		r := po.Lp
		return int(math.Ceil(float64(po.price) * po.Lp)), r, true
	}

	return 0, 0.0, false
}

func (po *Position) buyBack(c CandleStick) (int, float64, bool) {
	// LC
	if float64(c.Start) >= po.lossCutPrice() {
		r := -(1.0 - float64(po.price)/float64(c.Start))
		return -(c.Start - po.price), r, true
	}

	// LC
	if float64(c.High) >= po.lossCutPrice() {
		r := -po.Lc
		return int(math.Ceil(float64(po.price) * po.Lc)), r, true
	}

	// PROFIT
	if float64(c.Start) <= po.profitPrice() {
		r := (float64(po.price)/float64(c.Start) - 1)
		return -(c.Start - po.price), r, true
	}

	// PROFIT
	if float64(c.Low) <= po.profitPrice() {
		r := po.Lp
		return int(math.Ceil(float64(po.price) * po.Lp)), r, true
	}

	return 0, 0.0, false
}

func (po Position) lossCutPrice() float64 {
	if po.t == buy {
		return float64(po.price) * (1 - po.Lc)
	}
	// sell
	return float64(po.price) * (1 + po.Lc)
}

func (po Position) profitPrice() float64 {
	if po.t == buy {
		return float64(po.price) * (1 + po.Lp)
	}
	// sell
	return float64(po.price) * (1 - po.Lp)
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
