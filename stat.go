package main

import (
	"math"
	"sync"
	"time"
)

type Stat struct {
	History        []*Result
	lock           *sync.Mutex
	historyKeepMax int
	Count          int
	nowIndex       int //for slice loop use
	LatencySum     time.Duration
	DroppedNum     int
}

type Result struct {
	Latency time.Duration
	Dropped int //1 for dropped,0 for non-dropped ,just for easy caculate
}

func (stat *Stat) AddResult(latency time.Duration, dropped bool) {
	var d int
	if dropped {
		d = 1
	} else {
		d = 0
	}
	stat.addResult(&Result{
		Latency: latency,
		Dropped: d,
	})
}

func (stat *Stat) addResult(newresult *Result) {
	stat.lock.Lock()
	defer stat.lock.Unlock()
	oldresult := stat.History[stat.nowIndex]
	var (
		oldlatency time.Duration
		olddropped int
	)
	if oldresult == nil {
		oldlatency, olddropped = 0, 0
	} else {
		oldlatency = oldresult.Latency
		olddropped = oldresult.Dropped
	}
	stat.History[stat.nowIndex] = newresult

	stat.LatencySum = stat.LatencySum - oldlatency + newresult.Latency
	stat.DroppedNum = stat.DroppedNum - olddropped + newresult.Dropped
	stat.nowIndex++
	if stat.nowIndex >= stat.historyKeepMax {
		stat.nowIndex = 0
	}
	stat.Count++
	if stat.Count >= stat.historyKeepMax {
		stat.Count = stat.historyKeepMax
	}

}

func NewStat(max int) *Stat {
	return &Stat{
		History:        make([]*Result, max),
		lock:           &sync.Mutex{},
		historyKeepMax: max,
		Count:          0,
		nowIndex:       0,
		LatencySum:     0,
		DroppedNum:     0,
	}
}

func (stat *Stat) DropRate() float64 {
	if stat.Count == 0 {
		return 0
	}
	return float64(stat.DroppedNum) / float64(stat.Count)
}

func (stat *Stat) AverageLatency() time.Duration {
	if stat.Count != 0 && stat.DroppedNum == stat.Count {
		return time.Second * 10
	}
	if stat.Count == stat.DroppedNum {
		return time.Nanosecond
	}
	return stat.LatencySum / time.Duration(stat.Count-stat.DroppedNum)
}

func (stat *Stat) Score() float64 {
	ping := float64(stat.AverageLatency())
	lost := float64(stat.DropRate())
	return math.Pow((1-lost), 50) / ping
}
