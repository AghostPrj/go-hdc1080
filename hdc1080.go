/**
 * @Author: aghost
 * @Author: ggg17226@gmail.com
 * @Date: 2021/12/30 8:59
 * @Desc:
 */

package go_hdc1080

import (
	"github.com/AghostPrj/go-i2c"
	"time"
)

const (
	Hdc1080RegTemperature     = 0x0
	Hdc1080RegHumidity        = 0x1
	Hdc1080RegConfigAndStatus = 0x2
	Hdc1080RegSerialPart1     = 0xfb
	Hdc1080RegSerialPart2     = 0xfc
	Hdc1080RegSerialPart3     = 0xfd
	Hdc1080RegManufacturerId  = 0xfe
	Hdc1080RegDeviceId        = 0xff

	Hdc1080CommandDelay = 20 * time.Millisecond
)

type HDC1080 struct {
	fp *i2c.I2C
}

// NewHdc1080 新建一个hdc1080的连接句柄
// NewHdc1080 open a new connection of hdc1080
func NewHdc1080(addr uint8, bus int) (*HDC1080, error) {
	c, err := i2c.NewI2C(addr, bus)
	if err != nil {
		return nil, err
	}
	hdc1080 := HDC1080{
		fp: c,
	}
	return &hdc1080, nil
}

// Close 关闭hdc1080连接
// Close hdc1080 connection.
func (h *HDC1080) Close() error {
	return h.fp.Close()
}

func (h *HDC1080) GetDeviceId() (uint16, error) {
	return h.fp.ReadRegU16BEWithDelay(Hdc1080RegDeviceId, Hdc1080CommandDelay)
}
