package models

import (
	"strconv"
	"sync"
	"time"
)

var (
	LogTypeBid   = "bid"
	LogTypeWin   = "win"
	LogTypeClick = "click"
)

type TimeLog interface {
	LogTime() time.Time
}

type BidResponse struct {
	Type string    `json:"type"`
	Id   string    `json:"id"`
	Time time.Time `json:"time"`

	Targeting  *Targeting `json:"targeting"`
	Publisher  string     `json:"publisher"`
	CampaignID int        `json:"campaign_id"`
	AdUrl      string     `json:"ad_url"`

	BidPrice float64 `json:"bid_price"`
}

func NewBidResponse() *BidResponse {
	now := time.Now()
	campaignID := GetCampaignID()

	return &BidResponse{
		Type: LogTypeBid,
		Id:   strconv.Itoa(int(now.UnixNano())),
		Time: now,

		Targeting: &Targeting{
			Geo:         GetGeoTargeting(),
			Device:      GetDeviceTargeting(),
			Demographic: GetDemographicTargeting(),
		},
		Publisher:  GetPublisherUrl(),
		CampaignID: campaignID,
		AdUrl:      GetAdUrl(campaignID),
		BidPrice:   GetBidPrice(),
	}
}

func (br *BidResponse) LogTime() time.Time {
	return br.Time
}

type Targeting struct {
	Geo         *TargetingGeo         `json:"geo"`
	Device      *TargetingDevice      `json:"device"`
	Demographic *TargetingDemographic `json:"demographic"`
}

type TargetingGeo struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	Zip     string `json:"zip"`
}

type TargetingDevice struct {
	Type string `json:"type"`
	Os   string `json:"os"`
}

type TargetingDemographic struct {
	Gender   string `json:"gender"`
	AgeRange string `json:"age_range"`
}

type WinNotice struct {
	Type     string    `json:"type"`
	BidId    string    `json:"bid_id"`
	BidTime  time.Time `json:"bid_time"`
	Time     time.Time `json:"time"`
	WinPrice float64   `json:"win_price"`
}

func NewWinNotice(br *BidResponse) *WinNotice {
	return &WinNotice{
		Type:     LogTypeWin,
		BidId:    br.Id,
		BidTime:  br.Time,
		Time:     br.Time.Add(GetWinNoticeTimeDiff()),
		WinPrice: GetSecondPrice(br.BidPrice),
	}
}

func (wn *WinNotice) LogTime() time.Time {
	return wn.Time
}

type Click struct {
	Type    string    `json:"type"`
	BidId   string    `json:"bid_id"`
	BidTime time.Time `json:"bid_time"`
	Time    time.Time `json:"time"`
}

func NewClick(br *BidResponse) *Click {
	return &Click{
		Type:    LogTypeClick,
		BidId:   br.Id,
		BidTime: br.Time,
		Time:    br.Time.Add(GetClickTimeDiff()),
	}
}

func (cl *Click) LogTime() time.Time {
	return cl.Time
}

type TimeLogHeap []TimeLog

func (a TimeLogHeap) Len() int {
	return len(a)
}

func (a TimeLogHeap) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a TimeLogHeap) Less(i, j int) bool {
	return a[i].LogTime().Sub(a[j].LogTime()) < 0
}

func (a *TimeLogHeap) Push(x interface{}) {
	*a = append(*a, x.(TimeLog))
}

func (a *TimeLogHeap) Pop() interface{} {
	old := *a
	n := len(old)
	x := old[n-1]
	*a = old[0 : n-1]
	return x
}

func (a *TimeLogHeap) Peak() TimeLog {
	l := *a
	if len(l) == 0 {
		return nil
	}
	return l[0]
}

type TimeLogHeapSync struct {
	TimeLogHeap
	sync.Mutex
}

type ChoiceString struct {
	Weight int
	Value  string
}
