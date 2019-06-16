package core

import (
	"errors"
	"fmt"
	"log/syslog"
	"sync"
)

// Train specifications.
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

// Track specifications.
type Track struct {
	ID          int
	Length      float64
	MaxVelocity float64
	Slope       float64
	BendRadius  float64
	Tunnel      bool
	// Source         int
	// Target         int
	// TrafficLightID int
	// PlatformID     int
}

// Sensors contains dynamic data collected by train's sensors.
type Sensors struct {
	//state         state
	Timestamp     int64 //Up to nanoseconds
	Setpoint      float64
	Position      float64
	Velocity      float64
	Acceleration  float64
	TractionForce float64
	BrakingForce  float64
	TractionPower float64
	BrakingPower  float64
	Mass          float64
	TrackIndex    int //Current track position in core.tracks slice
	TrackID       int
	RelPosition   float64 //Relative to the current track
	Slope         float64
	BendRadius    float64
	Tunnel        bool
	Resistance    float64 //Basic + line resistance
	BasicRes      float64
	SlopeRes      float64
	CurveRes      float64
	TunnelRes     float64
	LineRes       float64 //slope + curve + tunnel resistance
	NumPassengers int
	/*TractionWork  float64
	BrakeWork     float64
	JerkRate      float64*/

	//NextSemaphoreSignal TrafficSignal //Next semaphore signal
}

type locks struct {
	train   sync.RWMutex
	tracks  sync.RWMutex
	sensors sync.RWMutex
}

// Core collects essential information for train automation
type Core struct {
	//status
	train   Train
	tracks  []Track
	sensors Sensors
	lock    *locks
	log     *syslog.Writer
	//alerts
}

// New declares and initialises a Core instance.
func New(log *syslog.Writer) (Core, error) {
	if log == nil {
		// Panic?
		return Core{}, errors.New("Attempt to declare a new Core. Log not provided (nil)")
	}

	log.Info("New Core initialised")
	return Core{
		lock: &locks{},
		log:  log,
	}, nil
}

// GetTrain returns Train specifications.
func (c *Core) GetTrain() (Train, error) {
	c.lock.train.RLock()
	train := c.train
	c.lock.train.RUnlock()

	//c.log.Info(fmt.Sprintf("Get Train (ID %d)", c.train.ID))
	return train, nil
}

// SetTrain sets new Train specifications.
func (c *Core) SetTrain(train Train) error {
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

// GetTrack returns Track specifications by its ID.
func (c *Core) GetTrack(position int) (Track, error) {
	c.lock.tracks.RLock()

	if c.tracks == nil {
		c.lock.tracks.RUnlock()
		fail := "Attempt to GetTrack. Core.tracks is (nil)"
		c.log.Warning(fail)
		return Track{}, errors.New(fail)
	}

	if position >= len(c.tracks) || position < 0 {
		c.lock.tracks.RUnlock()
		fail := fmt.Sprintf("Attempt to GetTrack. Position %d out of bounds", position)
		c.log.Warning(fail)
		return Track{}, errors.New(fail)
	}

	track := c.tracks[position]
	c.lock.tracks.RUnlock()

	//c.log.Info(fmt.Sprintf("Get Track (ID %d)", track.ID))
	return track, nil

}

// AddTracks adds an ordered slice of Track to drive through.
func (c *Core) AddTracks(tracks ...Track) error {
	/*Validate tracks
	if !validated {
		return errors.New("blablabla")
	}*/

	c.lock.tracks.Lock()
	c.tracks = append(c.tracks, tracks...)
	c.lock.tracks.Unlock()

	c.log.Info(fmt.Sprintf("Add %+v", tracks))
	return nil
}

// SetTracks sets an ordered slice of Track to drive through.
// Existing stored tracks are replaced by new ones
func (c *Core) SetTracks(tracks ...Track) error {
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

// DeleteTracks drops Tracks in memory.
func (c *Core) DeleteTracks() {
	c.lock.tracks.Lock()
	c.tracks = nil
	c.lock.tracks.Unlock()

	c.log.Info("Delete Tracks")
}

// GetSensors retrieve current sensors data.
func (c *Core) GetSensors() (Sensors, error) {
	c.lock.sensors.RLock()
	sensors := c.sensors
	c.lock.sensors.RUnlock()

	return sensors, nil
}

// SetSensors sets new sensor values.
func (c *Core) SetSensors(sensors Sensors) error {
	c.lock.sensors.Lock()
	c.sensors = sensors
	c.lock.sensors.Unlock()

	return nil
}
