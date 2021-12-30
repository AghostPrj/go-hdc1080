/**
 * @Author: aghost
 * @Author: ggg17226@gmail.com
 * @Date: 2021/12/30 8:59
 * @Desc:
 */

package go_hdc1080

import (
	"fmt"
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

	Hdc1080ModeNormal = 0b0
	Hdc1080ModeReset  = 0b1

	Hdc1080HeaterEnable  = 0b1
	Hdc1080HeaterDisable = 0b0

	Hdc1080AcquisitionSingle = 0b0
	Hdc1080AcquisitionBoth   = 0b1

	Hdc1080BatteryStatusHigh = 0b0
	Hdc1080BatteryStatusLow  = 0b1

	Hdc1080TemperatureResolution14Bit = 0b0
	Hdc1080TemperatureResolution11Bit = 0b1

	Hdc1080HumidityResolution14Bit = 0b0
	Hdc1080HumidityResolution11Bit = 0b1
	Hdc1080HumidityResolution8Bit  = 0b10
)

type HDC1080 struct {
	fp *i2c.I2C
}

// HDC1080Config hdc1080配置结构体
// HDC1080Config config struct
type HDC1080Config struct {
	Reset                 uint8 `json:"reset"`
	Heater                uint8 `json:"heater"`
	AcquisitionMode       uint8 `json:"acquisition_mode"`
	BatteryStatus         uint8 `json:"battery_status"`
	TemperatureResolution uint8 `json:"temperature_resolution"`
	HumidityResolution    uint8 `json:"humidity_resolution"`
}

// Marshal 根据配置结构体生成配置
// Marshal convert struct to code
func (conf *HDC1080Config) Marshal() uint16 {
	var result uint16
	result = 0

	if conf.Reset <= 1 {
		result |= uint16(conf.Reset) << 15
	}

	if conf.Heater <= 1 {
		result |= uint16(conf.Heater) << 13
	}

	if conf.AcquisitionMode <= 1 {
		result |= uint16(conf.AcquisitionMode) << 12
	}

	if conf.BatteryStatus <= 1 {
		result |= uint16(conf.BatteryStatus) << 11
	}

	if conf.TemperatureResolution <= 1 {
		result |= uint16(conf.TemperatureResolution) << 10
	}

	if conf.HumidityResolution == Hdc1080HumidityResolution14Bit ||
		conf.HumidityResolution == Hdc1080HumidityResolution11Bit ||
		conf.HumidityResolution == Hdc1080HumidityResolution8Bit {
		result &= 0b1111110000000000
		result |= uint16(conf.HumidityResolution) << 8
	}

	return result
}

func parseHdc1080Config(confByte uint16) *HDC1080Config {
	result := HDC1080Config{}

	result.Reset = uint8(confByte >> 15 & Hdc1080ModeReset)
	result.Heater = uint8(confByte >> 13 & Hdc1080HeaterEnable)
	result.AcquisitionMode = uint8(confByte >> 12 & Hdc1080AcquisitionBoth)
	result.BatteryStatus = uint8(confByte >> 11 & Hdc1080BatteryStatusLow)
	result.TemperatureResolution = uint8(confByte >> 10 & Hdc1080TemperatureResolution11Bit)
	tmpHres := uint8(confByte >> 8 & 0b11)
	switch tmpHres {
	case Hdc1080HumidityResolution8Bit:
		result.HumidityResolution = Hdc1080HumidityResolution8Bit
		break
	case Hdc1080HumidityResolution11Bit:
		result.HumidityResolution = Hdc1080HumidityResolution11Bit
		break
	case Hdc1080HumidityResolution14Bit:
		result.HumidityResolution = Hdc1080HumidityResolution14Bit
		break
	default:
		break
	}

	return &result
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

// GetDeviceId 获取器件id，固定为0x1050
// GetDeviceId get device id, result is 0x1050
func (h *HDC1080) GetDeviceId() (uint16, error) {
	return h.fp.ReadRegU16BEWithDelay(Hdc1080RegDeviceId, Hdc1080CommandDelay)
}

// GetManufacturerId 获取制造商id，固定为0x5449
// GetManufacturerId get manufacturer id, result is 0x5449
func (h *HDC1080) GetManufacturerId() (uint16, error) {
	return h.fp.ReadRegU16BEWithDelay(Hdc1080RegManufacturerId, Hdc1080CommandDelay)
}

// GetSerialId 获取器件唯一id
// GetSerialId get device serial id
func (h *HDC1080) GetSerialId() (string, error) {
	serial1, err := h.fp.ReadRegU16BEWithDelay(Hdc1080RegSerialPart1, Hdc1080CommandDelay)
	if err != nil {
		return "", err
	}
	serial2, err := h.fp.ReadRegU16BEWithDelay(Hdc1080RegSerialPart2, Hdc1080CommandDelay)
	if err != nil {
		return "", err
	}
	serial3, err := h.fp.ReadRegU16BEWithDelay(Hdc1080RegSerialPart3, Hdc1080CommandDelay)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("%04X%04X%02X", serial1, serial2, serial3>>8)
	return result, nil
}

// GetConfig 获取配置
// GetConfig get chip config
func (h *HDC1080) GetConfig() (*HDC1080Config, error) {
	configBytes, err := h.fp.ReadRegU16BEWithDelay(Hdc1080RegConfigAndStatus, Hdc1080CommandDelay)
	if err != nil {
		return nil, err
	}
	return parseHdc1080Config(configBytes), nil
}

// SetConfig 写入配置
// SetConfig write config to chip
func (h *HDC1080) SetConfig(conf *HDC1080Config) error {
	confBytes := conf.Marshal()
	return h.fp.WriteRegU16BE(Hdc1080RegConfigAndStatus, confBytes)
}

// Reset 重置芯片配置
// Reset reset chip config
func (h *HDC1080) Reset() error {
	resetConf := HDC1080Config{Reset: Hdc1080ModeReset}
	return h.SetConfig(&resetConf)
}

// GetTemperature 获取摄氏温度
// GetTemperature get celsius temperature
func (h *HDC1080) GetTemperature() (float64, error) {
	tempBytes, err := h.fp.ReadRegU16BEWithDelay(Hdc1080RegTemperature, Hdc1080CommandDelay)
	if err != nil {
		return 0, err
	}
	result := float64(tempBytes)*160/65536 - 40
	return result, nil
}

// GetHumidity 获取相对湿度
// GetHumidity get relative humidity
func (h *HDC1080) GetHumidity() (float64, error) {
	tempBytes, err := h.fp.ReadRegU16BEWithDelay(Hdc1080RegHumidity, Hdc1080CommandDelay)
	if err != nil {
		return 0, err
	}
	result := float64(tempBytes) * 100 / 65536
	return result, nil
}
