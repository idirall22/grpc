package sample

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/idirall22/grpc/pb"
)

// NewKeyboard create new keyboard
func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}

// NewCPU create a cpu
func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	cores := randomInt(2, 8)
	threads := randomInt(cores, 12)
	minGhz := randomFloat64(2.2, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	return &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(cores),
		NumberThreads: uint32(threads),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}
}

// NewGPU create a new gpu
func NewGPU() *pb.GPU {
	brand := randomCPUBrand()
	name := randomGPUName(brand)
	minGhz := randomFloat64(1.0, 1.2)
	maxGhz := randomFloat64(minGhz, 2.0)
	memory := &pb.Memory{
		Value: uint64(randomInt(2, 8)),
		Unit:  pb.Memory_GIGABYTE,
	}

	return &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}
}

// NewRAM create a memory ram
func NewRAM() *pb.Memory {
	return &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_GIGABYTE,
	}
}

// NewSSd create a ssd storage
func NewSSd() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

// NewHDD create a hdd storage
func NewHDD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

// NewScreen create a new screen
func NewScreen() *pb.Screen {
	height := randomInt(1080, 4320)
	width := height * 16 / 9
	return &pb.Screen{
		Resolution: &pb.Screen_Resolution{
			Height: uint32(height),
			Width:  uint32(width),
		},
		SizeInch:   randomFloat32(7, 13),
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
}

// NewLaptop create a new laptop
func NewLaptop() *pb.Laptop {
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	return &pb.Laptop{
		Id:        randomID(),
		Brand:     brand,
		Name:      name,
		Cpu:       NewCPU(),
		Keyboard:  NewKeyboard(),
		Ram:       NewRAM(),
		Screen:    NewScreen(),
		Gpus:      []*pb.GPU{NewGPU()},
		Storages:  []*pb.Storage{NewHDD(), NewSSd()},
		Weight:    &pb.Laptop_WeightKg{WeightKg: 1},
		PriceUsd:  randomFloat64(1500, 3000),
		UpdatedAt: ptypes.TimestampNow(),
	}
}
