package atp

import (
	"math"
	"time"

	"github.com/trencat/Trencat/train/core"
)

const gravity float64 = 9.80665

// refresh updates real time data. The setpoint refers to acceleration.
func (atp *ATP) refresh() error {
	atp.lock.setpoint.RLock()
	setpoint := atp.setpoint
	atp.lock.setpoint.RUnlock()

	co := &(atp.core)
	prev, _ := co.GetSensors()
	train, _ := co.GetTrain()
	track, _ := co.GetTrack(prev.TrackIndex) // TODO Handle error!

	var new core.Sensors
	before := prev.Timestamp
	now := time.Now().UnixNano()
	deltaSec := float64(now-before) * 1e-9

	mass := train.Mass + float64(prev.NumPassengers)*70 //Add a bit of mass for each passenger

	new.Velocity = math.Max(0.0, prev.Velocity+deltaSec*prev.Acceleration)
	new.Position = prev.Position + 0.5*(prev.Velocity+new.Velocity)*deltaSec
	if setpoint.Value < 0.0 && math.Abs(new.Velocity) < 0.01 {
		new.Acceleration = 0
	} else {
		new.Acceleration = setpoint.Value //check limits!
	}
	new.BasicRes = mass * (train.ResistanceLin + train.ResistanceQua*new.Velocity*new.Velocity)
	new.SlopeRes = mass * gravity * math.Sin(track.Slope)

	if track.BendRadius <= 100 {
		//Prompt an Alert here? Danger!
	} else if track.BendRadius < 300 {
		new.CurveRes = 4.91 * mass / (track.BendRadius - 55)
	} else {
		new.CurveRes = 6.3 * mass / (track.BendRadius - 55)
	}

	if track.Tunnel {
		new.TunnelRes = 1.296 * 1e-9 * math.Max(track.Length-new.RelPosition, 0.0) * gravity * new.Velocity * new.Velocity
	} else {
		new.TunnelRes = 0
	}

	new.LineRes = new.SlopeRes + new.CurveRes + new.TunnelRes
	new.Resistance = new.BasicRes + new.LineRes

	force := mass*train.MassFactor*new.Acceleration + new.Resistance

	if force >= 0 {
		new.TractionForce = force
		new.TractionPower = new.TractionForce * new.Velocity
		new.BrakingForce = 0
		new.BrakingPower = 0
	} else {
		new.TractionForce = 0
		new.TractionPower = 0
		new.BrakingForce = -force
		new.BrakingPower = new.BrakingForce * new.Velocity
	}

	new.Setpoint = setpoint.Value
	new.Mass = mass
	new.Slope = track.Slope
	new.BendRadius = track.BendRadius
	new.Tunnel = track.Tunnel
	new.Timestamp = now

	new.RelPosition = prev.RelPosition + 0.5*(prev.Velocity+new.Velocity)
	if new.RelPosition > track.Length {
		new.RelPosition = 0
		new.TrackIndex = prev.TrackIndex + 1
		nextTrack, _ := co.GetTrack(prev.TrackIndex + 1) //TODO: Handle error here
		new.TrackID = nextTrack.ID
	} else {
		new.TrackIndex = prev.TrackIndex
		new.TrackID = prev.TrackID
	}

	// Update
	atp.core.SetSensors(new)

	return nil
}

/*type RealTime struct {
	//Status? (Driving, stopped, alarm, etc)
	Timestamp       int64

	TractionWork        float64
	BrakeWork           float64
	JerkRate            float64
	NumPassengers       int
	NextSemaphoreSignal TrafficSignal //Next semaphore signal
}*/
