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
	NextTrackID int
	PrevTrackID int
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
	tracks  map[int]Track
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
		tracks: make(map[int]Track),
		lock:   &locks{},
		log:    log,
	}, nil
}

// GetTrain returns Train specifications.
func (c *Core) GetTrain() (Train, error) {
	c.lock.train.RLock()
	defer c.lock.train.RUnlock()

	//c.log.Info(fmt.Sprintf("Get Train (ID %d)", c.train.ID))
	return c.train, nil
}

// SetTrain sets new Train specifications.
func (c *Core) SetTrain(train Train) error {
	/*Validate train
	if !validated {
		return errors.New("blablabla")
	}*/

	c.lock.train.Lock()
	defer c.lock.train.Unlock()

	c.train = train
	c.log.Info(fmt.Sprintf("Set Train%+v", c.train))
	return nil
}

// GetTrack returns Track specifications by its ID.
func (c *Core) GetTrack(ID int) (Track, error) {
	c.lock.tracks.RLock()
	defer c.lock.tracks.RUnlock()

	if c.tracks == nil {
		fail := "Attempt to GetTrack, but Core.tracks is not initialised (nil)"
		c.log.Warning(fail)
		return Track{}, errors.New(fail)
	}

	track, exists := c.tracks[ID]
	if !exists {
		fail := fmt.Sprintf("Attempt to get Track %d. ID doesn't exist", ID)
		c.log.Warning(fail)
		return Track{}, errors.New(fail)
	}

	//c.log.Info(fmt.Sprintf("Get Track (ID %d)", track.ID))
	return track, nil

}

// InsertTrack sets Track specifications.
func (c *Core) InsertTrack(track Track) error {
	/*Validate track
	if !validated {
		return errors.New("blablabla")
	}*/

	c.lock.tracks.Lock()
	defer c.lock.tracks.Unlock()

	if c.tracks == nil {
		fail := fmt.Sprintf("Attempt to insert Track%+v. Core.tracks is not initialised (nil). ", track)
		c.log.Warning(fail)
		return errors.New(fail)
	}

	if prevTrack, exists := c.tracks[track.PrevTrackID]; exists {
		if prevTrack.NextTrackID != track.ID {
			fail := fmt.Sprintf("Attempt to insert Track%+v. Doesn't match with existing Track%+v.", track, prevTrack)
			c.log.Warning(fail)
			return errors.New(fail)
		}
	}
	if nextTrack, exists := c.tracks[track.NextTrackID]; exists {
		if nextTrack.PrevTrackID != track.ID {
			fail := fmt.Sprintf("Attempt to insert Track%+v. Doesn't match with existing Track%+v.", track, nextTrack)
			c.log.Warning(fail)
			return errors.New(fail)
		}
	}

	c.log.Info(fmt.Sprintf("Insert Track%+v", track))
	c.tracks[track.ID] = track
	return nil
}

// DeleteTrack drops Track given its ID.
func (c *Core) DeleteTrack(ID int) {
	c.lock.tracks.Lock()
	defer c.lock.tracks.RUnlock()

	if _, exists := c.tracks[ID]; !exists {
		fail := fmt.Sprintf("Attempt to delete track %d. ID doesn't exist", ID)
		c.log.Warning(fail)
	}

	c.log.Info(fmt.Sprintf("Delete Track (ID %d)", ID))
	delete(c.tracks, ID)
}

// GetSensors retrieve current sensors data.
func (c *Core) GetSensors() (Sensors, error) {
	c.lock.sensors.RLock()
	defer c.lock.sensors.RUnlock()

	return c.sensors, nil
}

// SetSensors sets new sensor values.
func (c *Core) SetSensors(sensors Sensors) error {
	c.lock.sensors.Lock()
	defer c.lock.sensors.Unlock()

	c.sensors = sensors
	return nil
}
