package atp

import (
	"errors"
	"fmt"
	"log/syslog"
	"sync"
	"time"

	"github.com/trencat/Trencat/train/core"
)

// Setpoint contains value and timestamp in nanosecods.
type Setpoint struct {
	Value     float64
	Timestamp int64
}

type communications struct {
	setpoint struct {
		channel          <-chan Setpoint
		stop             chan struct{}
		stopNotification chan<- struct{}
	}
	sensorsChannels map[int]sensorsChannel
}

type sensorsChannel struct {
	ID        int
	frequency int
	channel   chan core.Sensors
	stop      chan struct{}
}

type locks struct {
	setpoint sync.RWMutex
	comms    struct {
		setpoint        sync.RWMutex
		sensorsChannels sync.RWMutex
	}
}

// ATP implements Automatic Train Protection
type ATP struct {
	core     core.Core
	setpoint Setpoint
	comms    communications
	lock     *locks
	log      *syslog.Writer
}

// New declares and initialises an ATP object.
func New(log *syslog.Writer) (ATP, error) {
	co, error := core.New(log)
	if error != nil {
		fail := fmt.Sprintf("Attempt to declare a new ATP. %s", error)
		//panic?
		return ATP{}, errors.New(fail)
	}

	newATP := ATP{
		core: co,
		lock: &locks{},
		log:  log,
	}
	newATP.comms.sensorsChannels = make(map[int]sensorsChannel)
	log.Info("New ATP initialised")
	return newATP, nil
}

// GetTrain returns Train specifications.
func (atp *ATP) GetTrain() (core.Train, error) {
	return atp.core.GetTrain()
}

// SetTrain sets new Train specifications.
func (atp *ATP) SetTrain(train core.Train) error {
	return atp.core.SetTrain(train)
}

// GetTrack returns Track specifications by its ID.
func (atp *ATP) GetTrack(ID int) (core.Track, error) {
	return atp.core.GetTrack(ID)
}

// InsertTrack sets Track specifications.
func (atp *ATP) InsertTrack(track core.Track) error {
	return atp.core.InsertTrack(track)
}

// DeleteTrack drops Track given its ID.
func (atp *ATP) DeleteTrack(ID int) {
	atp.core.DeleteTrack(ID)
}

// OpenSetpointChannel creates a channel to deliver setpoints while driving.
// Call StopSetpointChannel to stop reading setpoints and free memory.
// The ATP can call atp.StopSetpointChannel whenever appropiate.
// TODO: Document examples
func (atp *ATP) OpenSetpointChannel() (chan<- Setpoint, <-chan struct{}, error) {
	atp.lock.comms.setpoint.Lock()

	if atp.comms.setpoint.channel != nil {
		atp.lock.comms.setpoint.Unlock()
		fail := "Attempt to open setpoint channel. Channel is already open"
		atp.log.Warning(fail)
		return nil, nil, errors.New(fail)
	}

	//TODO: Implement canOpen()
	// if error := c.canOpen(); error != nil {
	// 	atp.lock.comms.setpoint.Unlock()
	// 	fail := fmt.Sprintf("Attempt to open setpoint channel. %s", error)
	// 	c.log.Warning(fail)
	// 	return nil, errors.New(fail)
	// }

	channel := make(chan Setpoint)
	stopNotification := make(chan struct{})
	stop := make(chan struct{})

	atp.comms.setpoint.channel = channel
	atp.comms.setpoint.stopNotification = stopNotification
	atp.comms.setpoint.stop = stop

	atp.lock.comms.setpoint.Unlock()

	if error := atp.readSetpoints(); error != nil {
		fail := fmt.Sprintf("Attempt to read setpoints from channel. %s", error)
		atp.log.Warning(fail)

		atp.lock.comms.setpoint.Lock()
		defer atp.lock.comms.setpoint.Unlock()

		atp.comms.setpoint.channel = nil
		atp.comms.setpoint.stopNotification = nil
		atp.comms.setpoint.stop = nil
		return nil, nil, error
	}

	atp.log.Info("Open Setpoint channel")
	return channel, stopNotification, nil
}

func (atp *ATP) readSetpoints() error {
	atp.lock.comms.setpoint.RLock()
	defer atp.lock.comms.setpoint.RUnlock()

	if atp.comms.setpoint.channel == nil || atp.comms.setpoint.stop == nil {
		fail := "Attempt to listen to setpoint channel. Channel not initialised"
		atp.log.Warning(fail)
		return errors.New(fail)
	}

	go func() {
	loop:
		for {
			select {
			case sp, isOpen := <-atp.comms.setpoint.channel:
				if isOpen {
					// Update setpoint
					atp.lock.setpoint.Lock()
					atp.setpoint = sp
					atp.lock.setpoint.Unlock()
				}
			case <-atp.comms.setpoint.stop:
				// Notify user that setpoint channel will not be read anymore
				atp.log.Info("Dropping setpoint channel gracefully")
				close(atp.comms.setpoint.stopNotification)
				break loop
			}
		}

		//Safe cleanup
		atp.lock.comms.setpoint.Lock()
		defer atp.lock.comms.setpoint.Unlock()
		atp.comms.setpoint.channel = nil
		atp.comms.setpoint.stop = nil
		atp.log.Info("Dropped setpoint channel")
	}()

	return nil
}

// StopSetpointChannel stops ATP from reading setpoints and frees resources.
// It does not close the setpoint channel to prevent the sender from writing on a closed channel (which implies panic!)
func (atp *ATP) StopSetpointChannel() error {
	atp.lock.comms.setpoint.Lock()
	defer atp.lock.comms.setpoint.Unlock()

	if atp.comms.setpoint.stop == nil {
		fail := "Attempt to stop setpoint channel. Channel is not open"
		atp.log.Warning(fail)
		return errors.New(fail)
	}

	atp.log.Info("Stopping setpoint channel signal sent")
	close(atp.comms.setpoint.stop)

	return nil
}

// NewSensorChannel creates a channel that delivers RealTime sensors data at the given frequency rate.
func (atp *ATP) NewSensorChannel(ID int, frequency int) (<-chan core.Sensors, error) {

	atp.lock.comms.sensorsChannels.Lock()

	if atp.comms.sensorsChannels == nil {
		atp.lock.comms.sensorsChannels.Unlock()
		fail := fmt.Sprintf("Attempt to get sensor channel (ID: %d, freq: %dms). "+
			"Core.comms.sensorsChannels is not initialised (nil). ", ID, frequency)
		atp.log.Warning(fail)
		return nil, errors.New(fail)
	}

	if _, exists := atp.comms.sensorsChannels[ID]; exists {
		atp.lock.comms.sensorsChannels.Unlock()
		fail := fmt.Sprintf("Attempt to get sensor channel (ID: %d, freq: %dms). ID already exists", ID, frequency)
		atp.log.Warning(fail)
		return nil, errors.New(fail)
	}

	channel := make(chan core.Sensors)
	stop := make(chan struct{})

	sensChan := sensorsChannel{
		ID:        ID,
		frequency: frequency,
		channel:   channel,
		stop:      stop,
	}

	atp.comms.sensorsChannels[ID] = sensChan
	atp.lock.comms.sensorsChannels.Unlock()

	if error := atp.startSensorChannel(&sensChan); error != nil {
		fail := fmt.Sprintf("Attempt to start sensor channel. %s", error)
		atp.log.Warning(fail)

		atp.lock.comms.sensorsChannels.Lock()
		delete(atp.comms.sensorsChannels, ID)
		atp.lock.comms.sensorsChannels.Unlock()

		return nil, error
	}

	atp.log.Info(fmt.Sprintf("New sensor channel (ID %d, freq %dms)", ID, frequency))
	return channel, nil
}

func (atp *ATP) startSensorChannel(sensChan *sensorsChannel) error {
	//TODO: Validate frequency. Should be greater than a certain threshold.

	atp.lock.comms.sensorsChannels.RLock()
	defer atp.lock.comms.sensorsChannels.RUnlock()

	if sensChan.channel == nil || sensChan.stop == nil {
		fail := fmt.Sprintf("Attempt to start sensor channel %+v. Channels not initialised", &sensChan)
		atp.log.Warning(fail)
		return errors.New(fail)
	}

	go func(atp *ATP, sensChan *sensorsChannel) {
		ticker := time.NewTicker(time.Duration(sensChan.frequency) * time.Millisecond)

	loop:
		for {
			select {
			case <-ticker.C:
				sensors, _ := atp.core.GetSensors()
				select {
				case sensChan.channel <- sensors:
					//Successful deliver
				default:
					//Cannot deliver measurements.
				}
			case <-sensChan.stop:
				//Close channel
				ticker.Stop()
				atp.log.Info(fmt.Sprintf("Closing sensor channel (ID %d) gracefully", sensChan.ID))
				close(sensChan.channel)
				break loop
			}
		}

		//Safe cleanup
		atp.lock.comms.sensorsChannels.Lock()
		defer atp.lock.comms.sensorsChannels.Unlock()
		delete(atp.comms.sensorsChannels, sensChan.ID)
		atp.log.Info(fmt.Sprintf("Closed sensor channel (ID %d)", sensChan.ID))
	}(atp, sensChan)

	return nil
}

// CloseSensorChannel closes the RealTime's channel matching the given ID.
func (atp *ATP) CloseSensorChannel(ID int) error {
	atp.lock.comms.sensorsChannels.Lock()
	defer atp.lock.comms.sensorsChannels.Unlock()

	if atp.comms.sensorsChannels == nil {
		fail := fmt.Sprintf("Attempt to close sensor channel from real time (ID %d). "+
			"Core.comms.sensorsChannels is not initialised (nil). ", ID)
		atp.log.Warning(fail)
		return errors.New(fail)
	}

	sensChan, exists := atp.comms.sensorsChannels[ID]
	if !exists {
		fail := fmt.Sprintf("Attempt to close sensor channel from real time (ID %d). ID doesn't exist", ID)
		atp.log.Warning(fail)
		return errors.New(fail)
	}

	atp.log.Info(fmt.Sprintf("Closing sensor channel (ID %d) signal sent", ID))
	close(sensChan.stop)
	return nil
}

func (atp *ATP) Start() (<-chan struct{}, error) {

	// Quick fix
	atp.core.SetSensors(core.Sensors{
		Timestamp: time.Now().UnixNano(),
		TrackID:   1,
	})

	ticker := time.NewTicker(time.Duration(200) * time.Millisecond)
	stop := make(chan struct{})

	go func() {
	loop:
		for {
			select {
			case <-ticker.C:
				atp.refresh()
			case <-stop:
				ticker.Stop()
				break loop
			}
		}
	}()

	return stop, nil
}
