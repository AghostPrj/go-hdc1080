/**
 * @Author: aghost
 * @Author: ggg17226@gmail.com
 * @Date: 2021/12/30 9:12
 * @Desc:
 */

package go_hdc1080

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func createConnection(t *testing.T) *HDC1080 {
	hdc1080, err := NewHdc1080(0x40, 0)
	if err != nil {
		t.Errorf("create connection error: %s", err)
	}
	return hdc1080
}

func TestHDC1080_GetDeviceId(t *testing.T) {
	hdc1080 := createConnection(t)
	defer hdc1080.Close()
	id, err := hdc1080.GetDeviceId()
	if err != nil {
		t.Errorf("get hdc1080 device id error: %s", err)
	} else {
		logrus.WithField("op", "test").Infof("deviceId: 0x%04X", id)
	}

}

func TestHDC1080_GetManufacturerId(t *testing.T) {
	hdc1080 := createConnection(t)
	defer hdc1080.Close()
	id, err := hdc1080.GetManufacturerId()
	if err != nil {
		t.Errorf("get hdc1080 manufacturer id error: %s", err)
	} else {
		logrus.WithField("op", "test").Infof("manufacturerId: 0x%04X", id)
	}
}

func TestHDC1080_GetSerialId(t *testing.T) {
	hdc1080 := createConnection(t)
	defer hdc1080.Close()
	id, err := hdc1080.GetSerialId()
	if err != nil {
		t.Errorf("get hdc1080 serial id error: %s", err)
	} else {
		logrus.WithField("op", "test").Infof("serialId: %s", id)
	}
}

func TestHDC1080_Config(t *testing.T) {
	hdc1080 := createConnection(t)
	defer hdc1080.Close()

	testConf := HDC1080Config{
		Reset:                 Hdc1080ModeNormal,
		Heater:                Hdc1080HeaterEnable,
		AcquisitionMode:       Hdc1080AcquisitionSingle,
		TemperatureResolution: Hdc1080TemperatureResolution11Bit,
		HumidityResolution:    Hdc1080HumidityResolution8Bit,
	}

	err := hdc1080.SetConfig(&testConf)
	if err != nil {
		t.Errorf("set hdc1080 config error: %s", err)
	}

	conf, err := hdc1080.GetConfig()
	marshal, _ := json.Marshal(conf)
	if err != nil {
		t.Errorf("get hdc1080 config error: %s", err)
	} else {
		logrus.WithField("op", "test").Infof("config: %s", marshal)
	}

	errReset := hdc1080.Reset()
	if errReset != nil {
		t.Errorf("set hdc1080 config error: %s", errReset)
	}
	time.Sleep(time.Millisecond * 50)
	conf, err = hdc1080.GetConfig()
	marshal, _ = json.Marshal(conf)
	if err != nil {
		t.Errorf("get hdc1080 config error: %s", err)
	} else {
		logrus.WithField("op", "test").Infof("config: %s", marshal)
	}
}

func TestHDC1080ConfigAndStatus_Marshal(t *testing.T) {
	config := HDC1080Config{AcquisitionMode: Hdc1080AcquisitionBoth}
	logrus.WithField("op", "test").Infof("config and status: 0x%04X", config.Marshal())
	hdc1080Config := parseHdc1080Config(0x1000)
	marshal, _ := json.Marshal(*hdc1080Config)
	logrus.WithField("op", "test").Info(string(marshal))
}

func TestHDC1080_GetTemperature(t *testing.T) {
	hdc1080 := createConnection(t)
	defer hdc1080.Close()
	temp, err := hdc1080.GetTemperature()
	if err != nil {
		t.Errorf("get hdc1080 temperature error: %s", err)
	}
	logrus.WithField("op", "test").Infof("temperature: %0.4f℃", temp)
}

func TestHDC1080_GetHumidity(t *testing.T) {
	hdc1080 := createConnection(t)
	defer hdc1080.Close()
	hum, err := hdc1080.GetHumidity()
	if err != nil {
		t.Errorf("get hdc1080 humidity error: %s", err)
	}
	logrus.WithField("op", "test").Infof("relative humidity: %0.4f", hum)
}
