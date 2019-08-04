package sms

import (
	"errors"
	"fmt"
	"github.com/andern/keysms"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/config"
	"vicinity-tinymesh-vas-co2/vas-co2-backend/model"
)

const (
	threshold = 800
)

type SMS struct {
	db     *gorm.DB
	config *config.SMSConfig
}

type Warning struct {
	UpdatedAt time.Time
}

func New(smsConfig *config.SMSConfig, db *gorm.DB) *SMS {
	return &SMS{
		db:     db,
		config: smsConfig,
	}
}

func empty(field string) bool {
	return field == ""
}

func (sms *SMS) Auth() error {
	user := sms.config.User
	key := sms.config.Key

	if empty(user) {
		return errors.New("keysms: username is not set")
	}

	if empty(key) {
		return errors.New("keysms: api key is not set")
	}

	keysms.Auth(user, key)
	return nil
}

func (sms *SMS) notify(eid string) error {

	message := fmt.Sprintf("CO2 %s value is above 800", eid)

	response, err := keysms.SendSMS(message, sms.config.Recipients...)
	if err != nil {
		return err
	}

	if !response.OK {
		return errors.New("notify key sms failed: " + response.Message.Message)
	}

	return nil
}

func (sms *SMS) ResetSend(oid uuid.UUID) error {
	var sensor model.Sensor

	sms.db.Where(model.Sensor{Oid: oid}).First(&sensor)

	if sensor.Oid != oid {
		return errors.New(oid.String() + " no such oid in database")
	}

	sensor.NotificationSent = false
	sms.db.Save(&sensor)

	return nil
}

func (sms *SMS) GetLastWarning(oid uuid.UUID) (bool, error) {
	return sms.isSent(oid)
}

func (sms *SMS) isSent(oid uuid.UUID) (bool, error) {
	var sensor model.Sensor
	sms.db.Where(model.Sensor{Oid: oid}).First(&sensor)

	if sensor.Oid != oid {
		return false, errors.New(oid.String() + " no such oid in database")
	}

	return sensor.NotificationSent, nil
}

func (sms *SMS) NotifyOnce(oid uuid.UUID, eid string, value int) error {

	if value < threshold {
		return nil
	}

	if sent, err := sms.isSent(oid); err != nil {
		return err
	} else if sent {
		// skip send if already sent
		log.Printf("skipping sms for oid '%s' \n", oid.String())
		return nil
	}

	return sms.notifyOnce(oid, eid)
}

func (sms *SMS) notifyOnce(oid uuid.UUID, eid string) error {
	var sensor model.Sensor

	// Save for KPI
	var alert = model.Notification{SensorOid: oid, TimeNotified: time.Now()}
	sms.db.Create(&alert)

	sms.db.Where(model.Sensor{Oid: oid}).First(&sensor)

	err := sms.notify(eid)
	if err != nil {
		return err
	}

	log.Printf("%s: %s sms sent to %v\n", oid.String(), eid, sms.config.Recipients)

	sensor.NotificationSent = true

	sms.db.Save(&sensor)
	return nil
}
