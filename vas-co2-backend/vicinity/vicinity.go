package vicinity

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/config"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/model"
)

type Client struct {
	config *config.VicinityConfig
	db     *gorm.DB
	td     *gin.H
}

type EventData struct {
	Value        int    `json:"value" binding:"required"`
	Unit         string `json:"unit"`
	Milliseconds string `json:"timestamp" binding:"required"`
}

type chartData struct {
	T     time.Time
	Value int
}

type dateRange struct {
	T time.Time
}

const (
	coreService = "core:Service"
	version     = "1.0.0"
)

type VAS struct {
	Oid        string        `json:"oid"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Version    string        `json:"version"`
	Keywords   []string      `json:"keywords"`
	Properties []interface{} `json:"properties"`
	Events     []interface{} `json:"events"`
	Actions    []interface{} `json:"actions"`
}

func (c *Client) makeVAS(oid, name, version string, kw []string) VAS {
	return VAS{
		Oid:      oid,
		Name:     name,
		Type:     coreService,
		Version:  version,
		Keywords: kw,
		// rest is empty
		Properties: []interface{}{},
		Events:     []interface{}{},
		Actions:    []interface{}{},
	}
}

func (c *Client) GetThingDescription() *gin.H {
	if c.td == nil {

		var vasGroup []VAS
		vasGroup = append(vasGroup, c.makeVAS(c.config.Oid, "TinyMesh VAS - CWi CO2 Monitor", version, []string{"co2"}))

		c.td = &gin.H{
			"adapter-id":         c.config.AdapterID,
			"thing-descriptions": vasGroup,
		}
	}

	return c.td
}

func (c *Client) GetDateRange(oid uuid.UUID) *gin.H {
	var result []dateRange
	c.db.Raw(
		`SELECT DATE_TRUNC('day', time) as t
		FROM readings r 
		WHERE r.sensor_oid = ?
		GROUP BY t
		ORDER BY t
		ASC`, oid).Scan(&result)

	var days []time.Time

	for _, day := range result {
		days = append(days, day.T)
	}

	return &gin.H{
		"days": days,
	}
}

func (c *Client) GetReadingsByDate(oid uuid.UUID, dateString string) (*gin.H, error) {
	var labels []string
	var data []int
	var result []chartData

	c.db.Raw(
		`select time as t, value
		FROM readings r
		WHERE r.time::date = ? 
		AND r.sensor_oid = ? 
		ORDER BY r.time ASC`, dateString, oid).Scan(&result)

	for _, row := range result {
		labels = append(labels, row.T.Format("15:04"))
		data = append(data, row.Value)
	}

	readings := &gin.H{
		"labels": labels,
		"data":   data,
	}

	return readings, nil
}

func (c *Client) GetReadings(oid uuid.UUID) (*gin.H, error) {
	var result []chartData
	var labels []string
	var data []int

	lowerTime := time.Now().Add(time.Duration(-12) * time.Hour)
	upperTime := time.Now()

	c.db.Raw(
		`SELECT DATE_TRUNC('hour', time) as t, 
		ROUND(AVG(value), 0) as value 
		FROM readings r 
		WHERE r.time BETWEEN ? AND ?
		AND r.sensor_oid = ? 
		GROUP BY t 
		ORDER BY t
		ASC`, lowerTime, upperTime, oid).Scan(&result)

	if c.db.Error != nil {
		log.Println(c.db.Error.Error())
		return nil, errors.New("could not execute select query")
	}

	for _, row := range result {
		labels = append(labels, row.T.Format("15:04"))
		data = append(data, row.Value)
	}

	readings := &gin.H{
		"labels": labels,
		"data":   data,
	}

	return readings, nil
}

func makeTimestamp(milliseconds int64) time.Time {
	return time.Unix(0, milliseconds*int64(time.Millisecond))
}

func (c *Client) StoreEventData(e EventData, oid uuid.UUID, eid string) error {
	var sensor model.Sensor

	i, err := strconv.ParseInt(e.Milliseconds, 10, 64)
	if err != nil {
		log.Println(err.Error())
	}

	tm := makeTimestamp(i)

	c.db.Where(model.Sensor{Oid: oid}).FirstOrCreate(&sensor, model.Sensor{Oid: oid, Eid: eid, Unit: "ppm"})

	if c.db.Error != nil {
		log.Println(c.db.Error.Error())
		return errors.New(fmt.Sprintf("could not fetch/create oid %v", oid.String()))
	}

	log.Println(tm.String())

	c.db.Create(&model.Reading{Value: e.Value, Time: tm, SensorOid: sensor.Oid})

	if c.db.Error != nil {
		log.Println(c.db.Error.Error())
		return errors.New(fmt.Sprintf("could not store event reading of oid: %v", oid.String()))
	}

	return nil
}

func New(vicinityConfig *config.VicinityConfig, db *gorm.DB) *Client {
	return &Client{
		config: vicinityConfig,
		db:     db,
		td:     nil,
	}
}
