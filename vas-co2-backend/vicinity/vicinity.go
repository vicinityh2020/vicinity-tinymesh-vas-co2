package vicinity

import (
	"github.com/gin-gonic/gin"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/config"
)

type Client struct {
	config *config.VicinityConfig
	db     interface{}
	td     *gin.H
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

func New(vicinityConfig *config.VicinityConfig, db interface{}) *Client {
	return &Client{
		config: vicinityConfig,
		db:     db,
		td:     nil,
	}
}
