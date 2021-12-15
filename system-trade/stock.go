package systemtrade

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type CandleStick struct {
	date  time.Time
	start int
	high  int
	low   int
	end   int
}

func (c CandleStick) Date() time.Time {
	return c.date
}

type CandleSticks []CandleStick

func (c CandleSticks) DMA(d int)(r float64) {
	size := len(c) -1
	for i:=0; i<d; i++ {
		println(c[size-i].end)
	}
	return 
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
			// 日付,始値,高値,安値,終値,出来高,終値,調整値
			d := strings.Replace(recode[0], "-", "/", -1)
			layout := "2006/01/02"
			dt, _ := time.Parse(layout, d)
			s, _ := strconv.Atoi(recode[1])
			h, _ := strconv.Atoi(recode[2])
			l, _ := strconv.Atoi(recode[3])
			e, _ := strconv.Atoi(recode[4])
			c := CandleStick{
				date:  dt,
				start: s,
				high:  h,
				low:   l,
				end:   e,
			}
			// fmt.Printf("%T: %v\n", recode, recode)
			rc = append(rc, c)
		}

		f.Close()
	}
	return rc, nil
}
