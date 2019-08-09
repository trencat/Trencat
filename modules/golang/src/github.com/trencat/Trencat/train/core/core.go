// Package core provides a simple (yet complete) implementation
// of train movement. It provides implementations of:
//   - interfaces.Core
//   - interfaces.Train
//   - interfaces.Track
//   - interfaces.Sensors
//   - interfaces.InitialConditions
//   - interfaces.Setpoint
package core

import (
	"fmt"
	"log/syslog"
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/trencat/Trencat/train/interfaces"
)

// Train specifications. Implements interfaces.Train
type Train struct {
	ID            int
	Mass          float64
	MassFactor    float64
	Length        float64
	MaxForce      float64
	MaxBrake      float64
	ResistanceLin float64
	ResistanceQua float64
}

// Track specifications. Implements interfaces.Track
type Track struct {
	ID          int
	Length      float64
	MaxVelocity float64
	Slope       float64
	BendRadius  float64
	Tunnel      bool
}

// Setpoint implements interfaces.Setpoint
type Setpoint struct {
	Value float64
	Time  time.Time
}

// When returns the time when the setpoint was recorded.
func (s Setpoint) When() time.Time {
	return s.Time
}

// Sensors contains dynamic data collected by train's sensors.
// All values are expressed in the International System of Units.
// Implements interfaces.Sensors
type Sensors struct {
	Time          time.Time
	Setpoint      float64
	Position      float64 // Relative to the beginning of the track
	Velocity      float64
	Acceleration  float64
	TractionForce float64
	BrakingForce  float64
	TractionPower float64
	BrakingPower  float64
	Mass          float64
	TrackIndex    int // Current track position in core.tracks slice
	TrackID       int
	RelPosition   float64 // Relative to the current track
	Slope         float64
	BendRadius    float64
	Tunnel        bool
	Resistance    float64 // Basic + line resistance
	BasicRes      float64
	SlopeRes      float64
	CurveRes      float64
	TunnelRes     float64
	LineRes       float64 // slope + curve + tunnel resistance
	NumPassengers int
}

// When returns the time when the sensors were recorded.
func (s Sensors) When() time.Time {
	return s.Time
}

const gravity float64 = 9.80665

type locks struct {
	train   sync.RWMutex
	tracks  sync.RWMutex
	sensors sync.RWMutex
}

// Core collects essential information for train automation and
// implements train movement. Implements interfaces.Core.
type Core struct {
	train   Train
	tracks  []Track
	sensors Sensors
	lock    *locks
	log     *syslog.Writer
}

// New declares and initialises a Core instance.
func New(log *syslog.Writer) (Core, error) {
	if log == nil {
		// Panic?
		err := errors.New("Attempt to declare a new Core. Log not provided (nil)")
		fmt.Printf("%+v", err)
		return Core{}, err
	}

	log.Info("New Core initialised")
	return Core{
		sensors: Sensors{Time: time.Now()},
		lock:    &locks{},
		log:     log,
	}, nil
}

// GetTrain returns Train specifications.
func (c *Core) GetTrain() (interfaces.Train, error) {
	return c.getTrain()
}

func (c *Core) getTrain() (Train, error) {
	c.lock.train.RLock()
	train := c.train
	c.lock.train.RUnlock()

	//c.log.Info(fmt.Sprintf("Get Train (ID %d)", c.train.ID))
	return train, nil
}

// SetTrain sets new Train specifications.
func (c *Core) SetTrain(train interfaces.Train) error {
	t, ok := train.(Train)
	if !ok {
		err := errors.New("Attempt to SetTrain. interface type is not core.Train")
		c.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	return c.setTrain(t)
}

func (c *Core) setTrain(train Train) error {
	/*Validate train
	if !validated {
		return errors.New("blablabla")
	}*/

	c.lock.train.Lock()
	c.train = train
	c.lock.train.Unlock()

	c.log.Info(fmt.Sprintf("Set Train%+v", c.train))
	return nil
}

// SetInitConditions sets initial conditions of the Train(i.e. position,
// velocity, acceleration, etc.). It must be called before atp.Start method.
func (c *Core) SetInitConditions(conditions interfaces.InitConditions) error {
	cond, ok := conditions.(Sensors)
	if !ok {
		err := errors.New("Attempt to SetInitConditions. interface type is not core.Sensors")
		c.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	return c.setInitConditions(cond)
}

func (c *Core) setInitConditions(conditions Sensors) error {
	err := c.setSensors(conditions)

	if err == nil {
		c.log.Info(fmt.Sprintf("Set initial conditions Sensors%+v", conditions))
	}

	return err
}

// GetTrack returns Track specifications by its ID.
func (c *Core) GetTrack(position int) (interfaces.Track, error) {
	return c.getTrack(position)
}

func (c *Core) getTrack(position int) (Track, error) {
	c.lock.tracks.RLock()

	if c.tracks == nil {
		c.lock.tracks.RUnlock()
		err := errors.New("Attempt to GetTrack. Core.tracks is (nil)")
		c.log.Warning(fmt.Sprintf("%+v", err))
		return Track{}, err
	}

	if position >= len(c.tracks) || position < 0 {
		c.lock.tracks.RUnlock()
		err := errors.Errorf("Attempt to GetTrack. Position %d out of bounds", position)
		c.log.Warning(fmt.Sprintf("%+v", err))
		return Track{}, err
	}

	track := c.tracks[position]
	c.lock.tracks.RUnlock()

	//c.log.Info(fmt.Sprintf("Get Track (ID %d)", track.ID))
	return track, nil

}

// SetTracks sets an ordered slice of Track to drive through.
// Existing stored tracks are replaced by new ones
func (c *Core) SetTracks(tracks ...interfaces.Track) error {

	newTracks := make([]Track, len(tracks))

	for i, track := range tracks {
		t, ok := track.(Track)
		if !ok {
			err := errors.New("Attempt to SetTracks. interface type is not core.Track")
			c.log.Warning(fmt.Sprintf("%+v", err))
			return err
		}

		newTracks[i] = t
	}

	return c.setTracks(newTracks...)
}

func (c *Core) setTracks(tracks ...Track) error {
	/*Validate tracks
	if !validated {
		return errors.New("blablabla")
	}*/

	c.lock.tracks.Lock()
	c.tracks = tracks
	c.lock.tracks.Unlock()

	c.log.Info(fmt.Sprintf("Set %+v", tracks))
	return nil
}

// GetSensors retrieve current sensors data.
func (c *Core) GetSensors() (interfaces.Sensors, error) {
	return c.getSensors()
}

func (c *Core) getSensors() (Sensors, error) {
	c.lock.sensors.RLock()
	sensors := c.sensors
	c.lock.sensors.RUnlock()

	return sensors, nil
}

// SetSensors sets new sensor values. Useful to set Core initial conditions.
func (c *Core) SetSensors(sensors interfaces.Sensors) error {
	s, ok := sensors.(Sensors)
	if !ok {
		err := errors.New("Attempt to SetSensors. interface type is not core.Sensors")
		c.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	return c.setSensors(s)
}

func (c *Core) setSensors(sensors Sensors) error {
	c.lock.sensors.Lock()
	c.sensors = sensors
	c.lock.sensors.Unlock()

	return nil
}

// UpdateSensors updates real time data after a given time duration.
// This implementation assumes that the setpoint refers to acceleration.
func (c *Core) UpdateSensors(setpoint interfaces.Setpoint, elapsed time.Duration) (interfaces.Sensors, error) {
	sp, ok := setpoint.(Setpoint)
	if !ok {
		err := errors.New("Attempt to UpdateSensors. interface type is not core.Setpoint")
		c.log.Warning(fmt.Sprintf("%+v", err))
		return Sensors{}, err
	}

	return c.updateSensorsAcceleration(sp.Value, elapsed)
}

func (c *Core) updateSensorsAcceleration(setpoint float64, elapsed time.Duration) (Sensors, error) {

	//TODO: REMOVE HARDCODED CONSTANTS

	prev, _ := c.getSensors()               // No errors implemented still
	train, _ := c.getTrain()                // No errors implemented still
	track, _ := c.getTrack(prev.TrackIndex) // TODO Handle error!

	// Update track
	beginNewTrack := (prev.RelPosition > track.Length)
	if beginNewTrack {
		track, _ = c.getTrack(prev.TrackIndex + 1) //TODO: Handle error here
	}

	var new Sensors

	// Time
	deltaSec := elapsed.Seconds()
	new.Time = prev.Time.Add(elapsed)
	// TrackID
	new.TrackID = track.ID
	// TrackIndex
	if beginNewTrack {
		new.TrackIndex = prev.TrackIndex + 1
	} else {
		new.TrackIndex = prev.TrackIndex
	}
	//Number of passengers
	new.NumPassengers = prev.NumPassengers
	// Mass (add average mass for each passenger)
	new.Mass = train.Mass + float64(new.NumPassengers)*70
	// Setpoint
	new.Setpoint = setpoint
	// Velocity
	new.Velocity = math.Max(0.0, prev.Velocity+deltaSec*prev.Acceleration)
	if new.Velocity >= track.MaxVelocity {
		c.log.Warning(fmt.Sprintf("Current velocity %fm/s exceeds maximum velocity %fm/s", new.Velocity, track.MaxVelocity))
		// Todo: Trigger state machine warnings.
	}
	// Position
	new.Position = prev.Position + 0.5*(prev.Velocity+new.Velocity)*deltaSec
	// Relative position
	if beginNewTrack {
		new.RelPosition = 0.5 * (prev.Velocity + new.Velocity)
	} else {
		new.RelPosition = prev.RelPosition + 0.5*(prev.Velocity+new.Velocity)
	}
	// Slope
	new.Slope = track.Slope
	// Bend Radius
	new.BendRadius = track.BendRadius
	// Tunnel
	new.Tunnel = track.Tunnel
	// Basic resistance
	new.BasicRes = new.Mass * (train.ResistanceLin + train.ResistanceQua*new.Velocity*new.Velocity)
	// Slope resistance
	new.SlopeRes = new.Mass * gravity * math.Sin(new.Slope)
	// Curve resistance
	if track.BendRadius <= 100 {
		//Prompt an Alert here? Danger!
	} else if track.BendRadius < 300 {
		new.CurveRes = 4.91 * new.Mass / (new.BendRadius - 55)
	} else {
		new.CurveRes = 6.3 * new.Mass / (new.BendRadius - 55)
	}
	// Tunnel resistance
	if track.Tunnel {
		new.TunnelRes = 1.296 * 1e-9 * math.Max(track.Length-new.RelPosition, 0.0) * gravity * new.Velocity * new.Velocity
	} else {
		new.TunnelRes = 0
	}
	// Line resistance
	new.LineRes = new.SlopeRes + new.CurveRes + new.TunnelRes
	new.Resistance = new.BasicRes + new.LineRes
	// Acceleration
	maxAcceleration := (train.MaxForce - new.Resistance) / (new.Mass * train.MassFactor)
	maxDeceleration := ((-1)*train.MaxBrake - new.Resistance) / (new.Mass * train.MassFactor)
	if setpoint > 0.0 && setpoint > maxAcceleration {
		// Correction required. Accelerating more than allowed
		c.log.Warning(fmt.Sprintf("Acceleration setpoint %fm/s2 exceeds maximum acceleration %fm/s", setpoint, maxAcceleration))
		// Todo: Trigger state machine warnings.
		new.Acceleration = maxAcceleration
	} else if setpoint < 0.0 && setpoint < maxDeceleration {
		// Correction required. Decelerating more than allowed
		c.log.Warning(fmt.Sprintf("Deceleration setpoint %fm/s2 exceeds maximum deceleration %fm/s", setpoint, maxDeceleration))
		// Todo: Trigger state machine warnings.
		new.Acceleration = maxDeceleration
	} else if setpoint < 0.0 && new.Velocity < 0.01 {
		// Reverse gear not allowed.
		new.Acceleration = 0
		new.Velocity = 0
	} else {
		// Setpoint within limits
		new.Acceleration = setpoint
	}
	// Force & power
	force := new.Mass*train.MassFactor*new.Acceleration + new.Resistance
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

	// Update
	c.setSensors(new)

	return new, nil
}
