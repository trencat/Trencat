package atp_test

import (
	"fmt"
	"log/syslog"
	"math"
	"os"
	"testing"
	"time"

	"github.com/trencat/Trencat/train/atp"
	"github.com/trencat/Trencat/train/core"
)

var log *syslog.Writer

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	// Setup
	syslog, error := syslog.Dial("tcp", "localhost:514",
		syslog.LOG_WARNING|syslog.LOG_LOCAL0, "coreTest")

	if error != nil {
		panic(fmt.Sprintf("%s", error))
	}

	log = syslog

	//Teardown
	os.Exit(m.Run())
}

func TestSimpleDriving(t *testing.T) {
	ATP, err := atp.New(log)
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = ATP.SetTrain(core.Train{
		ID:            1,
		Length:        75,
		Mass:          5.07e5,
		MassFactor:    1.06,
		MaxForce:      3e5,
		MaxBrake:      4.475e5,
		ResistanceLin: 0.014 / 5.07e5,
		ResistanceQua: 2.564e-5 / 5.07e5,
	})

	if err != nil {
		t.Fatalf("%s", err)
	}

	err = ATP.SetTracks(core.Track{
		ID:          1,
		Length:      10000,
		MaxVelocity: 14,
		Slope:       0,
		BendRadius:  math.Inf(1),
		Tunnel:      false,
	})

	if err != nil {
		t.Fatalf("%s", err)
	}

	_, err = ATP.Start()

	if err != nil {
		t.Fatalf("%s", err)
	}

	setpoint, _, _ := ATP.OpenSetpointChannel()

	log.Info("Starting in 1 seconds")
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
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	// Get current sensor readings
	sensorsChan, err := ATP.NewSensorChannel(1, 700)
	if err != nil {
		t.Fatalf("%s", err)
	}
	sensor := <-sensorsChan
	err = ATP.CloseSensorChannel(1)
	if err != nil {
		t.Errorf("%s", err)
	}

	log.Info(fmt.Sprintf("P:%.7f\tV:%.7f\tA:%.7f\tTf:%.7f\tBf:%.7f\tRes:%.7f\tTime:%f\n", sensor.Position, sensor.Velocity, sensor.Acceleration, sensor.TractionForce, sensor.BrakingForce, sensor.Resistance, float64(sensor.Timestamp)*1e-9))

	if math.Abs(sensor.Position-150.0) > 5e-1 {
		t.Errorf("Expect train to stop at position 150.0. Finally stopped at %f", sensor.Position)
	}

	if math.Abs(sensor.Velocity-0.0) > 1e-1 {
		t.Errorf("Expect train to have velocity 0.0. Velocity is %f", sensor.Velocity)
	}

	if math.Abs(sensor.Acceleration-0.0) > 1e-1 {
		t.Errorf("Expect train to have acceleration 0.0. Acceleration is %f", sensor.Acceleration)
	}

}
