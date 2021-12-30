/**
 * @Author: aghost
 * @Author: ggg17226@gmail.com
 * @Date: 2021/12/30 9:12
 * @Desc:
 */

package go_hdc1080

import (
	"github.com/sirupsen/logrus"
	"testing"
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
		logrus.WithField("op", "test").Infof("deviceId: %0X", id)
	}

}
