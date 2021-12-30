# go-hdc1080

德州仪器HDC1080温度湿度传感器golang库  
ti hdc1080 temperature and humidity sensor golang library


----------------------------
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

这个库使用 [golang](https://golang.org/) 编写，用于操作hdc1080温湿度传感器。

This library written in [Go programming language](https://golang.org/) intended to operation hdc1080 sensor chip.

-------------------

## Usage

```go
func main(){
hdc1080, err := NewHdc1080(0x40, 0)
if err != nil {
// ....
}
defer hdc1080.Close()
temp, err := hdc1080.GetTemperature()
hum, err := hdc1080.GetHumidity()
}
```

详细使用方法参考 [hdc1080_test.go](./hdc1080_test.go)  




