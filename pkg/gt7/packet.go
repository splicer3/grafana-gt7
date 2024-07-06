package gt7

import (
	"encoding/binary"
	"encoding/json"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"golang.org/x/crypto/salsa20"
	"math"
	"time"
)

type TelemetryFrame struct {
	PackageID         int32
	BestLap           int32
	LastLap           int32
	CurrentLap        int16
	CurrentGear       uint8
	SuggestedGear     uint8
	FuelCapacity      float32
	CurrentFuel       float32
	Boost             float32
	TyreDiameterFL    float32
	TyreDiameterFR    float32
	TyreDiameterRL    float32
	TyreDiameterRR    float32
	TyreSpeedFL       float32
	TyreSpeedFR       float32
	TyreSpeedRL       float32
	TyreSpeedRR       float32
	CarSpeed          float32
	TyreSlipRatioFL   float32
	TyreSlipRatioFR   float32
	TyreSlipRatioRL   float32
	TyreSlipRatioRR   float32
	TimeOnTrack       time.Duration
	TotalLaps         int16
	CurrentPosition   int16
	TotalPositions    int16
	CarID             int32
	Throttle          float32
	RPM               float32
	RPMRevWarning     uint16
	Brake             float32
	RPMRevLimiter     uint16
	EstimatedTopSpeed int16
	Clutch            float32
	ClutchEngaged     float32
	RPMAfterClutch    float32
	OilTemp           float32
	WaterTemp         float32
	OilPressure       float32
	RideHeight        float32
	TyreTempFL        float32
	TyreTempFR        float32
	TyreTempRL        float32
	TyreTempRR        float32
	SuspensionFL      float32
	SuspensionFR      float32
	SuspensionRL      float32
	SuspensionRR      float32
	Gear1             float32
	Gear2             float32
	Gear3             float32
	Gear4             float32
	Gear5             float32
	Gear6             float32
	Gear7             float32
	Gear8             float32
	PositionX         float32
	PositionY         float32
	PositionZ         float32
	VelocityX         float32
	VelocityY         float32
	VelocityZ         float32
	RotationPitch     float32
	RotationYaw       float32
	RotationRoll      float32
	AngularVelocityX  float32
	AngularVelocityY  float32
	AngularVelocityZ  float32
	IsPaused          bool
	InRace            bool
}

func salsa20Dec(dat []byte) []byte {
	keyStr := "Simulator Interface Packet GT7 ver 0.0"
	var key [32]byte
	copy(key[:], keyStr[:32])

	// Seed IV is always located here
	oiv := dat[0x40:0x44]
	iv1 := binary.LittleEndian.Uint32(oiv)

	// Notice DEADBEAF, not DEADBEEF
	iv2 := iv1 ^ 0xDEADBEAF

	iv := make([]byte, 8)
	binary.LittleEndian.PutUint32(iv[0:4], iv2)
	binary.LittleEndian.PutUint32(iv[4:8], iv1)

	ddata := make([]byte, len(dat))
	salsa20.XORKeyStream(ddata, dat, iv, &key)

	magic := binary.LittleEndian.Uint32(ddata[0:4])
	if magic != 0x47375330 {
		log.DefaultLogger.Warn("Magic number mismatch")
		return []byte{}
	}

	return ddata
}

func ReadPacket(b []byte) (*TelemetryFrame, error) {
	dFrame := salsa20Dec(b)

	frame := convertTelemetryValues(dFrame)

	return frame, nil
}

func convertTelemetryValues(data []byte) *TelemetryFrame {
	returnedFrame := TelemetryFrame{
		PackageID:         int32(binary.LittleEndian.Uint32(data[0x70:0x74])),
		BestLap:           int32(binary.LittleEndian.Uint32(data[0x78:0x7C])),
		LastLap:           int32(binary.LittleEndian.Uint32(data[0x7C:0x80])),
		CurrentLap:        int16(binary.LittleEndian.Uint16(data[0x74:0x76])),
		CurrentGear:       data[0x90] & 0b00001111,
		SuggestedGear:     data[0x90] >> 4,
		FuelCapacity:      math.Float32frombits(binary.LittleEndian.Uint32(data[0x48:0x4C])),
		CurrentFuel:       math.Float32frombits(binary.LittleEndian.Uint32(data[0x44:0x48])),
		Boost:             math.Float32frombits(binary.LittleEndian.Uint32(data[0x50:0x54])) - 1,
		TyreDiameterFL:    math.Float32frombits(binary.LittleEndian.Uint32(data[0xB4:0xB8])),
		TyreDiameterFR:    math.Float32frombits(binary.LittleEndian.Uint32(data[0xB8:0xBC])),
		TyreDiameterRL:    math.Float32frombits(binary.LittleEndian.Uint32(data[0xBC:0xC0])),
		TyreDiameterRR:    math.Float32frombits(binary.LittleEndian.Uint32(data[0xC0:0xC4])),
		CarSpeed:          3.6 * math.Float32frombits(binary.LittleEndian.Uint32(data[0x4C:0x50])),
		TotalLaps:         int16(binary.LittleEndian.Uint16(data[0x76:0x78])),
		CurrentPosition:   int16(binary.LittleEndian.Uint16(data[0x84:0x86])),
		TotalPositions:    int16(binary.LittleEndian.Uint16(data[0x86:0x88])),
		CarID:             int32(binary.LittleEndian.Uint32(data[0x124:0x128])),
		Throttle:          float32(data[0x91]) / 2.55,
		RPM:               math.Float32frombits(binary.LittleEndian.Uint32(data[0x3C:0x40])),
		RPMRevWarning:     binary.LittleEndian.Uint16(data[0x88:0x8A]),
		Brake:             float32(data[0x92]) / 2.55,
		RPMRevLimiter:     binary.LittleEndian.Uint16(data[0x8A:0x8C]),
		EstimatedTopSpeed: int16(binary.LittleEndian.Uint16(data[0x8C:0x8E])),
		OilTemp:           math.Float32frombits(binary.LittleEndian.Uint32(data[0x5C:0x60])),
		WaterTemp:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x58:0x5C])),
		OilPressure:       math.Float32frombits(binary.LittleEndian.Uint32(data[0x54:0x58])),
		RideHeight:        1000 * math.Float32frombits(binary.LittleEndian.Uint32(data[0x38:0x3C])),
		TyreTempFL:        math.Float32frombits(binary.LittleEndian.Uint32(data[0x60:0x64])),
		TyreTempFR:        math.Float32frombits(binary.LittleEndian.Uint32(data[0x64:0x68])),
		TyreTempRL:        math.Float32frombits(binary.LittleEndian.Uint32(data[0x68:0x6C])),
		TyreTempRR:        math.Float32frombits(binary.LittleEndian.Uint32(data[0x6C:0x70])),
		SuspensionFL:      math.Float32frombits(binary.LittleEndian.Uint32(data[0xC4:0xC8])),
		SuspensionFR:      math.Float32frombits(binary.LittleEndian.Uint32(data[0xC8:0xCC])),
		SuspensionRL:      math.Float32frombits(binary.LittleEndian.Uint32(data[0xCC:0xD0])),
		SuspensionRR:      math.Float32frombits(binary.LittleEndian.Uint32(data[0xD0:0xD4])),
		PositionX:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x04:0x08])),
		PositionY:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x08:0x0C])),
		PositionZ:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x0C:0x10])),
		VelocityX:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x10:0x14])),
		VelocityY:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x14:0x18])),
		VelocityZ:         math.Float32frombits(binary.LittleEndian.Uint32(data[0x18:0x1C])),
		RotationPitch:     math.Float32frombits(binary.LittleEndian.Uint32(data[0x1C:0x20])),
		RotationYaw:       math.Float32frombits(binary.LittleEndian.Uint32(data[0x20:0x24])),
		RotationRoll:      math.Float32frombits(binary.LittleEndian.Uint32(data[0x24:0x28])),
		AngularVelocityX:  math.Float32frombits(binary.LittleEndian.Uint32(data[0x2C:0x30])),
		AngularVelocityY:  math.Float32frombits(binary.LittleEndian.Uint32(data[0x30:0x34])),
		AngularVelocityZ:  math.Float32frombits(binary.LittleEndian.Uint32(data[0x34:0x38])),
		IsPaused:          (data[0x8E] & 0b00000010) != 0,
		InRace:            (data[0x8E] & 0b00000001) != 0,
	}

	// Tyre speed calculation
	returnedFrame.TyreSpeedFL = float32(math.Abs(float64(3.6 * returnedFrame.TyreDiameterFL * math.Float32frombits(binary.LittleEndian.Uint32(data[0xA4:0xA8])))))
	returnedFrame.TyreSpeedFR = float32(math.Abs(float64(3.6 * returnedFrame.TyreDiameterFR * math.Float32frombits(binary.LittleEndian.Uint32(data[0xA8:0xAC])))))
	returnedFrame.TyreSpeedRL = float32(math.Abs(float64(3.6 * returnedFrame.TyreDiameterRL * math.Float32frombits(binary.LittleEndian.Uint32(data[0xAC:0xB0])))))
	returnedFrame.TyreSpeedRR = float32(math.Abs(float64(3.6 * returnedFrame.TyreDiameterRR * math.Float32frombits(binary.LittleEndian.Uint32(data[0xB0:0xB4])))))

	// Tyre slip ratio calculation
	if returnedFrame.CarSpeed > 0 {
		returnedFrame.TyreSlipRatioFL = returnedFrame.TyreSpeedFL / returnedFrame.CarSpeed
		returnedFrame.TyreSlipRatioFR = returnedFrame.TyreSpeedFR / returnedFrame.CarSpeed
		returnedFrame.TyreSlipRatioRL = returnedFrame.TyreSpeedRL / returnedFrame.CarSpeed
		returnedFrame.TyreSlipRatioRR = returnedFrame.TyreSpeedRR / returnedFrame.CarSpeed
	} else {
		returnedFrame.TyreSlipRatioFL = -1 // Using -1 to represent '  –  '
		returnedFrame.TyreSlipRatioFR = -1
		returnedFrame.TyreSlipRatioRL = -1
		returnedFrame.TyreSlipRatioRR = -1
	}

	return &returnedFrame
}

func telemetryFrameToMap(frame TelemetryFrame) map[string]float32 {
	var frameMap map[string]float32
	frameJson, err := json.Marshal(&frame)
	if err != nil {
		log.DefaultLogger.Error("Error converting frame", "error", err)
	}
	json.Unmarshal(frameJson, &frameMap)
	return frameMap
}

func TelemetryToDataFrame(tf TelemetryFrame) *data.Frame {
	frame := data.NewFrame("response")
	telemetryMap := telemetryFrameToMap(tf)

	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{time.Now()}),
	)

	for name, value := range telemetryMap {
		frame.Fields = append(frame.Fields,
			data.NewField(name, nil, []float32{value}),
		)
	}

	return frame
}
