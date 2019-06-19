package main

import (
	"fmt"
	"log"
	"log/syslog"
	"math"
	"time"

	"github.com/trencat/Trencat/train/atp"
	"github.com/trencat/Trencat/train/core"
)

func main() {
	syslog, err := syslog.Dial("tcp", "127.0.0.1:514",
		syslog.LOG_WARNING|syslog.LOG_LOCAL0, "ATP")
	if err != nil {
		log.Fatal(err)
	}

	train := core.Train{
		ID:            1,
		Length:        75,
		Mass:          5.07e5,
		MassFactor:    1.06,
		MaxForce:      3e5,
		MaxBrake:      4.475e5,
		ResistanceLin: 0.014 / 5.07e5,
		ResistanceQua: 2.564e-5 / 5.07e5,
	}

	track := core.Track{
		ID:          1,
		Length:      10000,
		MaxVelocity: 14,
		Slope:       0,
		BendRadius:  math.Inf(1),
		Tunnel:      false,
		// Source:      1,
		// Target:      2,
		// TrafficLightId: 1,
		// PlatformId:     0,
	}

	ATP, _ := atp.New(syslog)
	ATP.SetTrain(train)
	ATP.SetTracks(track)
	ATP.Start()

	go func() {
		sensorsChan, _ := ATP.NewSensorChannel(1, 700)
		for s := range sensorsChan {
			syslog.Info(fmt.Sprintf("P:%.7f\tV:%.7f\tA:%.7f\tTf:%.7f\tBf:%.7f\tRes:%.7f\tTime:%f\n", s.Position, s.Velocity, s.Acceleration, s.TractionForce, s.BrakingForce, s.Resistance, float64(s.Timestamp)*1e-9))
		}
	}()

	setpoint, _, _ := ATP.OpenSetpointChannel()
	syslog.Info("Starting in 1 seconds")
	time.Sleep(time.Duration(1) * time.Second)

	//Accelerate 5 seconds, cruise 10 seconds, brake 7
	for i := 0; i < 25; i++ {
		now := time.Now().UnixNano()
		if i < 5 { //2
			setpoint <- atp.Setpoint{Value: float64(2), Timestamp: now}
		} else if i < 15 {
			setpoint <- atp.Setpoint{Value: float64(0), Timestamp: now}
		} else if i < 20 {
			setpoint <- atp.Setpoint{Value: float64(-2), Timestamp: now}
		} else if i > 21 {
			ATP.StopSetpointChannel()
			err = ATP.CloseSensorChannel(1)
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}
