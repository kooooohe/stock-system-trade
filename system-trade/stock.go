package systemtrade

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	"time"
)

type CandleStick struct {
	Date  time.Time
	Start float64
	High  float64
	Low   float64
	End   float64
}

type CandleSticks []CandleStick

func (c CandleSticks) sort() CandleSticks {
	sort.Slice(c, func(i, j int) bool {
		return c[i].Date.Before(c[j].Date)
	})
	return c
}
func (c CandleSticks) DMA(d int, to int) (r float64) {
	sum := 0.0
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

func (r recode) StartDate() time.Time {
	return r.start.c.Date
}

func (r *recode) SetEnd(c CandleStick, p float64, d float64) {
	r.end = end{p: p, c: c}
	r.defference = d
}

type start struct {
	c CandleStick
	p float64
}

type end struct {
	c CandleStick
	p float64
}

type Score struct {
	win       int
	lose      int
	sum       float64
	buy       int
	shortSell int
	recodes   []recode
}

func (s *Score) SetStartRecode(c CandleStick, p float64, t positionType) {
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
	s.recodes[len(s.recodes)-1].SetEnd(c, 0, d)
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
	sum := 0.0
	sumS := 0.0
	sumB := 0.0
	cntB := 0
	cntS := 0

	winB := 0
	winS := 0
	loseB := 0
	loseS := 0

	tDate := s.recodes[0].StartDate()
	tYear := tDate.Year()

	for _, v := range s.recodes {
		if v.StartDate().Year() != tYear {
			fmt.Printf("【YAER】: %v\n", tYear)
			fmt.Printf("win:  %v (buy: %v shortsell: %v)\n", winB+winS, winB, winS)
			fmt.Printf("lose: %v (buy: %v shortsell: %v)\n", loseB+loseS, loseB, loseS)
			fmt.Printf("sum: %v\n\n", sum)
			fmt.Printf("sum Buy: %v\n", sumB)
			fmt.Printf("sumShortSell: %v\n\n", sumS)
			fmt.Printf("count: Buy: %v\n", cntB)
			fmt.Printf("count: ShortSell: %v\n", cntS)
			tYear = v.StartDate().Year()
			sum = 0
			winB = 0
			winS = 0
			loseB = 0
			loseS = 0

			sumS = 0
			sumB = 0
			cntB = 0
			cntS = 0
		}

		sum += v.defference

		isWin := v.defference > 0
		if isWin {
			if v.t == buy {
				winB++
			} else {
				winS++
			}
		} else {
			if v.t == buy {
				loseB++
			} else {
				loseS++
			}
		}

		if v.t == buy {
			cntB++
			sumB += v.defference
		} else {
			cntS++
			sumS += v.defference
		}

	}

	fmt.Printf("【YAER】: %v\n", s.recodes[len(s.recodes)-1].StartDate().Year())
	fmt.Printf("win: %v (buy: %v shortsell: %v)\n", winB+winS, winB, winS)
	fmt.Printf("lose: %v (buy: %v shortsell: %v)\n", loseB+loseS, loseB, loseS)
	fmt.Printf("sum: %v\n\n", sum)

	fmt.Printf("sum Buy: %v\n", sumB)
	fmt.Printf("sumShortSell: %v\n\n", sumS)
	fmt.Printf("count: Buy: %v\n", cntB)
	fmt.Printf("count: ShortSell: %v\n", cntS)

	/*TOTAL*/
	fmt.Println("==================")
	fmt.Printf("TOTAL sum: %v\n", s.sum)
	fmt.Printf("TOTAL win: %v\n", s.win)
	fmt.Printf("TOTAL lose: %v\n", s.lose)
	fmt.Printf("TOTAL buy: %v\n", s.buy)
	fmt.Printf("TOTAL shortSell: %v\n", s.shortSell)
	fmt.Printf("TOTAL len: %v\n", len(s.recodes))

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
	price float64
	Lc    float64
	Lp    float64 //指値
}

func (po *Position) Buy(p float64, c CandleStick) {
	po.t = buy
	po.price = p

	Result.SetStartRecode(c, p, buy)
}

func (po *Position) ShortSell(p float64, c CandleStick) {
	po.t = sell
	po.price = p

	Result.SetStartRecode(c, p, sell)
}

func (po *Position) Sell(c CandleStick) (float64, float64, bool) {
	a, b, ok := po.sell(c)
	if ok {
		po.t = nothing
		Result.SetEndRcode(c, b)
		if b > 0 {
			Result.Win()
		} else {

			Result.Lose()
		}
	}
	return a, b, ok
}

func (po *Position) BuyBack(c CandleStick) (float64, float64, bool) {
	a, b, ok := po.buyBack(c)
	if ok {
		Result.SetEndRcode(c, b)
		po.t = nothing
		if b > 0 {
			Result.Win()
		} else {

			Result.Lose()
		}
	}
	return a, b, ok
}
func (po *Position) sell(c CandleStick) (float64, float64, bool) {

	// LC
	if float64(c.Start) <= po.lossCutPrice() {
		r := -(1.0 - float64(c.Start)/float64(po.price))
		return c.Start - po.price, r, true
	}

	// LC
	if float64(c.Low) <= po.lossCutPrice() {
		r := -po.Lc
		return math.Ceil(float64(po.price) * po.Lc), r, true
	}

	// PROFIT
	if float64(c.Start) >= po.profitPrice() {
		r := (float64(c.Start)/float64(po.price) - 1)
		return c.Start - po.price, r, true
	}

	// PROFIT
	if float64(c.High) >= po.profitPrice() {
		r := po.Lp
		return math.Ceil(float64(po.price) * po.Lp), r, true
	}

	return 0, 0.0, false
}

func (po *Position) buyBack(c CandleStick) (float64, float64, bool) {
	// LC
	if float64(c.Start) >= po.lossCutPrice() {
		r := -(1.0 - float64(po.price)/float64(c.Start))
		return -(c.Start - po.price), r, true
	}

	// LC
	if float64(c.High) >= po.lossCutPrice() {
		r := -po.Lc
		return math.Ceil(float64(po.price) * po.Lc), r, true
	}

	// PROFIT
	if float64(c.Start) <= po.profitPrice() {
		r := (float64(po.price)/float64(c.Start) - 1)
		return -(c.Start - po.price), r, true
	}

	// PROFIT
	if float64(c.Low) <= po.profitPrice() {
		r := po.Lp
		return math.Ceil(float64(po.price) * po.Lp), r, true
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
func (po Position) Price() float64 {
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
			s, _ := strconv.ParseFloat(recode[1], 64)
			h, _ := strconv.ParseFloat(recode[2], 64)
			l, _ := strconv.ParseFloat(recode[3], 64)
			e, _ := strconv.ParseFloat(recode[4], 64)
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
