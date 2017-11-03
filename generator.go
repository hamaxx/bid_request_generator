package main

import (
	"container/heap"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hamaxx/bid_request_generator/models"
)

const (
	winRate   = 0.1
	clickRate = 0.01
)

func sendLog(w io.Writer, d interface{}) {
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		log.Fatal(err)
	}
}

func processTimeHeap(w io.Writer, h *models.TimeLogHeapSync) {
	for range time.Tick(time.Millisecond) {
		now := time.Now()

		x := h.Peak()
		if x == nil {
			continue
		}
		if x.LogTime().Sub(now) > 0 {
			continue
		}

		h.Lock()
		xi := heap.Pop(h)
		h.Unlock()

		sendLog(w, xi)
	}
}

func getRate() (int, int, error) {
	if len(os.Args) < 2 {
		return 1, 1, nil
	}

	rate, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return 0, 0, err
	}

	if len(os.Args) < 3 {
		return rate, 1, nil
	}

	proc, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return 0, 0, err
	}

	return rate, proc, nil
}

func runGenerator(w io.Writer, rate int, idx int) {
	winNoticeHeap := &models.TimeLogHeapSync{}
	go processTimeHeap(w, winNoticeHeap)

	clickHeap := &models.TimeLogHeapSync{}
	go processTimeHeap(w, clickHeap)

	count := 0
	t0 := time.Now()

	for range time.Tick(time.Second / time.Duration(rate)) {
		now := time.Now()
		if now.Sub(t0) > time.Second*10 {
			log.Printf("Real rate core %d: %f.2/sec", idx, float64(count)/now.Sub(t0).Seconds())
			t0 = now
			count = 0
		}
		count++

		br := models.NewBidResponse()
		sendLog(w, br)

		if rand.Float64() > winRate {
			continue
		}

		wn := models.NewWinNotice(br)

		winNoticeHeap.Lock()
		heap.Push(winNoticeHeap, wn)
		winNoticeHeap.Unlock()

		if rand.Float64() > clickRate {
			continue
		}

		cl := models.NewClick(br)

		clickHeap.Lock()
		heap.Push(clickHeap, cl)
		clickHeap.Unlock()
	}
}

func main() {
	w := os.Stdout

	rate, proc, err := getRate()
	if err != nil {
		log.Panicf("Input error: %s", err)
	}

	for i := 0; i < proc; i++ {
		go runGenerator(w, rate/proc, i)
	}

	select {}
}
