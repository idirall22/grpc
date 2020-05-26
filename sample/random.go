package sample

import (
	"math/rand"

	"github.com/google/uuid"

	"github.com/idirall22/grpc/pb"
)

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 0:
		return pb.Keyboard_UNKNOWN
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(s ...string) string {
	l := len(s)
	if l == 0 {
		return ""
	}
	return s[rand.Intn(l)]
}

func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"core i9 9900k",
			"core i7 9900k",
			"core i5 9900k",
			"core i3 9900k",
		)
	}
	return randomStringFromSet(
		"Ryzen 7 Pro 2200U",
		"Ryzen 5 Pro 3200U",
		"Ryzen 3 Pro 3200U",
	)
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD")
}

func randomGPUName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFromSet(
			"RTX 2080",
			"RTX 2070",
			"GTX 1660 Ti",
			"GTX 1070 Ti",
		)
	}
	return randomStringFromSet(
		"RX 590",
		"RX 580",
		"RX 5700-XT",
		"RX VEGA-56",
	)
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}
func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "Xps")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1")
	}
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomID() string {
	return uuid.New().String()
}
