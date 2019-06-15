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
		NextTrackID: 2,
		PrevTrackID: 0,
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
	ATP.InsertTrack(track)
	ATP.Start()

	setpoint, _, _ := ATP.OpenSetpointChannel()
	sensorsChan, _ := ATP.NewSensorChannel(1, 1500)

	//Accelerate 5 seconds, cruise 10 seconds, brake 7
	ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)

	var count int

loop:
	for {
		select {
		case rt := <-sensorsChan:
			syslog.Info(fmt.Sprintf("P:%.7f\tV:%.7f\tA:%.7f\tTf:%.7f\tBf:%.7f\tRes:%.7f\tTime:%f\n", rt.Position, rt.Velocity, rt.Acceleration, rt.TractionForce, rt.BrakingForce, rt.Resistance, float64(rt.Timestamp)*1e-9))

		case <-ticker.C:
			now := time.Now().UnixNano()
			if count < 5 { //2
				setpoint <- atp.Setpoint{Value: float64(2), Timestamp: now}
			} else if count < 15 {
				setpoint <- atp.Setpoint{Value: float64(0), Timestamp: now}
			} else if count < 20 {
				setpoint <- atp.Setpoint{Value: float64(-2), Timestamp: now}
			} else if count > 23 {
				ticker.Stop()
				rt := <-sensorsChan
				syslog.Info(fmt.Sprintf("P:%.7f\tV:%.7f\tA:%.7f\tTf:%.7f\tBf:%.7f\tRes:%.7f\tTime:%f\n", rt.Position, rt.Velocity, rt.Acceleration, rt.TractionForce, rt.BrakingForce, rt.Resistance, float64(rt.Timestamp)*1e-9))
				ATP.StopSetpointChannel()
				ATP.CloseSensorChannel(1)
				break loop
			}
			count++
		}
	}

	time.Sleep(2 * time.Second)

}
