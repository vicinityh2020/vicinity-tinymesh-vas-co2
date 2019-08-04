package vicinity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/config"
)

const (
	timeout = 5 * time.Second
	base    = "https://cpsgw.cs.uni-kl.de"
	pilot   = "Oslo"

	rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"
)

type CountWrapper struct {
	Amount int
}

type KPITracker struct {
	db         *gorm.DB
	config     *config.VicinityConfig
	httpClient *http.Client
	quit       chan bool
	log		   *log.Logger
}

type DatValue struct {
	Type   string `json:"type,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

func NewKPITracker(vicinityConfig *config.VicinityConfig, db *gorm.DB, logger io.Writer) *KPITracker {
	return &KPITracker{
		config: vicinityConfig,
		db:     db,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		log: log.New(logger, "", log.Ldate|log.Ltime),
		quit: make(chan bool),
	}
}

// goroutine
func (tracker *KPITracker) Tick(every time.Duration) {
	ticker := time.NewTicker(every * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				tracker.log.Println("Sending KPI ...")
				tracker.GatherAndReport()
			case <-tracker.quit:
				ticker.Stop()
				tracker.log.Println("Tracker stopped ...")
				return
			}
		}
	}()
}

func (tracker *KPITracker) Stop() {
	tracker.quit <- true
}

func (tracker *KPITracker) GatherAndReport() {
	t := time.Now()
	past := t.Add(-24 * time.Hour)

	var notifications CountWrapper
	tracker.db.Raw(`SELECT COUNT(*) as amount FROM notifications WHERE time_notified BETWEEN ? AND ?`, past, t).Scan(&notifications)

	var messages CountWrapper
	tracker.db.Raw(`SELECT COUNT(*) as amount FROM readings r WHERE r.time BETWEEN ? AND ?`, past, t).Scan(&messages)

	numOrgEtc := []DatValue{
		{Type: "Organizations", Amount: 1},
		{Type: "Devices", Amount: 26},
		{Type: "VAS", Amount: 5},
		{Type: "Friendships", Amount: 4},
		{Type: "Contracts", Amount: 1},
	}

	connectedDevices := &DatValue{Type: "device", Amount: 6}
	participants := &DatValue{Type: "Cleaning", Amount: 1}

	// Hardcoded KPIs
	tracker.Report(&t, "ConnectedDevices", connectedDevices)
	for _, v := range numOrgEtc {
		tracker.Report(&t, "NumOrgEtc", v)
	}
	tracker.Report(&t, "ParticipantsAmount", participants)
	tracker.Report(&t, "NumMaintenanceAlerts", 0)

	// Dynamically gathered KPIs
	tracker.Report(&t, "NumMsgReceived", messages.Amount)
	tracker.Report(&t, "NumNotification", notifications.Amount)
}

func (tracker *KPITracker) Report(t *time.Time, graphID string, num interface{}) {
	url := fmt.Sprintf("%s/vicinity-dashboard/Data/%s", base, pilot)

	body := map[string]interface{}{
		"timestamp": t.Format(rfc2822),
		"graph_id":  graphID,
		"dat_value": num,
	}

	payload, err := json.Marshal(&body)

	if err != nil {
		tracker.log.Println("Could not marshal in ..." + graphID)
		tracker.log.Println(err.Error())
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))

	if err != nil {
		tracker.log.Println("Could not create request in..." + graphID)
		tracker.log.Println(err.Error())
		return
	}

	req.Header.Add("AUTH", tracker.config.KPIKey)

	res, err := tracker.httpClient.Do(req)

	if err != nil {
		tracker.log.Println(err.Error())
		return
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		tracker.log.Print("Tracker", graphID, res.StatusCode, "\n")
	}

	text, err := ioutil.ReadAll(res.Body)

	if err != nil {
		tracker.log.Println(err.Error())
		return
	}

	tracker.log.Println(graphID + ": " + string(text))
}
