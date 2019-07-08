// Package atp provides a security layer over train movement.
// It implements interfaces.ATP.
package atp

import (
	"fmt"
	"log/syslog"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/trencat/Trencat/train/core"
	"github.com/trencat/Trencat/train/interfaces"
)

// communications gathers interfaces.setpoint and
// interfaces.Sensors channels
type communications struct {
	setpoint struct {
		channel          <-chan interfaces.Setpoint
		stop             chan struct{}
		stopNotification chan<- struct{}
	}
	sensorsChannels map[int]sensorsChannel
}

type sensorsChannel struct {
	ID        int
	frequency int
	channel   chan interfaces.Sensors
	stop      chan struct{}
}

type locks struct {
	setpoint sync.RWMutex
	comms    struct {
		setpoint        sync.RWMutex
		sensorsChannels sync.RWMutex
	}
}

// Atp implements interfaces.ATP.
type Atp struct {
	core     interfaces.Core
	setpoint interfaces.Setpoint
	comms    communications
	lock     *locks
	log      *syslog.Writer
}

// New declares and initialises an ATP object.
func New(log *syslog.Writer) (Atp, error) {
	co, err := core.New(log)
	if err != nil {
		//panic?
		return Atp{}, err
	}

	newATP := Atp{
		core:     &co,
		setpoint: core.Setpoint{Value: 0.0, Time: time.Unix(0, 0)},
		lock:     &locks{},
		log:      log,
	}
	newATP.comms.sensorsChannels = make(map[int]sensorsChannel)

	log.Info("New ATP initialised")
	return newATP, nil
}

// GetTrain returns Train specifications.
func (atp *Atp) GetTrain() (interfaces.Train, error) {
	return atp.core.GetTrain()
}

// SetTrain sets new Train specifications.
func (atp *Atp) SetTrain(train interfaces.Train) error {
	return atp.core.SetTrain(train)
}

// GetTrack returns Track specifications by its ID.
func (atp *Atp) GetTrack(position int) (interfaces.Track, error) {
	return atp.core.GetTrack(position)
}

// SetTracks sets an ordered slice of Track to drive through.
func (atp *Atp) SetTracks(track ...interfaces.Track) error {
	return atp.core.SetTracks(track...)
}

// SetInitConditions sets initial conditions of the Train(i.e. position,
// velocity, acceleration, etc.). It must be called before atp.Start method.
func (atp *Atp) SetInitConditions(conditions interfaces.InitConditions) error {
	return atp.core.SetInitConditions(conditions)
}

// getSetpoint returns last setpoint received in setpoint atp.comms.setpoint.channel.
func (atp *Atp) getSetpoint() interfaces.Setpoint {
	atp.lock.setpoint.RLock()
	setpoint := atp.setpoint
	atp.lock.setpoint.RUnlock()

	return setpoint
}

// OpenSetpointChannel creates a channel to deliver setpoints while driving.
// Call StopSetpointChannel to stop reading setpoints and free memory.
// The ATP can call atp.StopSetpointChannel whenever appropiate.
// TODO: Document examples
func (atp *Atp) OpenSetpointChannel() (chan<- interfaces.Setpoint, <-chan struct{}, error) {
	atp.lock.comms.setpoint.Lock()

	if atp.comms.setpoint.channel != nil {
		atp.lock.comms.setpoint.Unlock()
		err := errors.New("Attempt to open setpoint channel. Channel is already open")
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return nil, nil, err
	}

	//TODO: Implement canOpen()
	// if error := c.canOpen(); error != nil {
	// 	atp.lock.comms.setpoint.Unlock()
	// 	fail := fmt.Sprintf("Attempt to open setpoint channel. %s", error)
	// 	c.log.Warning(fail)
	// 	return nil, errors.New(fail)
	// }

	channel := make(chan interfaces.Setpoint)
	stopNotification := make(chan struct{})
	stop := make(chan struct{})

	atp.comms.setpoint.channel = channel
	atp.comms.setpoint.stopNotification = stopNotification
	atp.comms.setpoint.stop = stop
	atp.lock.comms.setpoint.Unlock()

	if err := atp.readSetpoints(); err != nil {
		atp.lock.comms.setpoint.Lock()
		atp.comms.setpoint.channel = nil
		atp.comms.setpoint.stopNotification = nil
		atp.comms.setpoint.stop = nil
		atp.lock.comms.setpoint.Unlock()
		return nil, nil, err
	}

	atp.log.Info("Open Setpoint channel")
	return channel, stopNotification, nil
}

// readSetpoints launches a new go routine to listen to atp.comms.setpoint.channel.
func (atp *Atp) readSetpoints() error {
	atp.lock.comms.setpoint.RLock()

	if atp.comms.setpoint.channel == nil || atp.comms.setpoint.stop == nil {
		atp.lock.comms.setpoint.Unlock()
		err := errors.New("Attempt to listen to setpoint channel. Channel not initialised")
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return err
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
		atp.comms.setpoint.channel = nil
		atp.comms.setpoint.stop = nil
		atp.lock.comms.setpoint.Unlock()
		atp.log.Info("Dropped setpoint channel")
	}()

	atp.lock.comms.setpoint.RUnlock()
	return nil
}

// StopSetpointChannel stops ATP from reading setpoints and frees resources.
// It does not close the setpoint channel to prevent the sender from writing on a closed channel (which implies panic!)
func (atp *Atp) StopSetpointChannel() error {
	atp.lock.comms.setpoint.Lock()

	if atp.comms.setpoint.stop == nil {
		atp.lock.comms.setpoint.Unlock()
		err := errors.New("Attempt to stop setpoint channel. Channel is not open")
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	close(atp.comms.setpoint.stop)
	atp.lock.comms.setpoint.Unlock()
	atp.log.Info("Stopping setpoint channel signal sent")

	return nil
}

// NewSensorChannel creates a channel that delivers RealTime sensors data at the given frequency rate.
func (atp *Atp) NewSensorChannel(ID int, frequency int) (<-chan interfaces.Sensors, error) {

	atp.lock.comms.sensorsChannels.Lock()

	if atp.comms.sensorsChannels == nil {
		atp.lock.comms.sensorsChannels.Unlock()
		err := errors.Errorf("Attempt to get sensor channel (ID: %d, freq: %dms). "+
			"Core.comms.sensorsChannels is not initialised (nil). ", ID, frequency)
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return nil, err
	}

	if _, exists := atp.comms.sensorsChannels[ID]; exists {
		atp.lock.comms.sensorsChannels.Unlock()
		err := errors.Errorf("Attempt to get sensor channel (ID: %d, freq: %dms). ID already exists", ID, frequency)
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return nil, err
	}

	channel := make(chan interfaces.Sensors)
	stop := make(chan struct{})

	sensChan := sensorsChannel{
		ID:        ID,
		frequency: frequency,
		channel:   channel,
		stop:      stop,
	}

	atp.comms.sensorsChannels[ID] = sensChan
	atp.lock.comms.sensorsChannels.Unlock()

	if err := atp.startSensorChannel(&sensChan); err != nil {
		atp.lock.comms.sensorsChannels.Lock()
		delete(atp.comms.sensorsChannels, ID)
		atp.lock.comms.sensorsChannels.Unlock()

		return nil, err
	}

	atp.log.Info(fmt.Sprintf("New sensor channel (ID %d, freq %dms)", ID, frequency))
	return channel, nil
}

func (atp *Atp) startSensorChannel(sensChan *sensorsChannel) error {
	//TODO: Validate frequency. Should be greater than a certain threshold.

	atp.lock.comms.sensorsChannels.RLock()

	if sensChan.channel == nil || sensChan.stop == nil {
		atp.lock.comms.sensorsChannels.RUnlock()
		err := errors.Errorf("Attempt to start sensorsChannel%+v. Channels not initialised", &sensChan)
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	go func(atp *Atp, sensChan *sensorsChannel) {
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
		delete(atp.comms.sensorsChannels, sensChan.ID)
		atp.lock.comms.sensorsChannels.Unlock()
		atp.log.Info(fmt.Sprintf("Closed sensor channel (ID %d)", sensChan.ID))
	}(atp, sensChan)

	atp.lock.comms.sensorsChannels.RUnlock()
	return nil
}

// CloseSensorChannel closes the RealTime's channel matching the given ID.
func (atp *Atp) CloseSensorChannel(ID int) error {
	atp.lock.comms.sensorsChannels.Lock()

	if atp.comms.sensorsChannels == nil {
		atp.lock.comms.sensorsChannels.Unlock()
		err := errors.Errorf("Attempt to close sensor channel from real time (ID %d). "+
			"Core.comms.sensorsChannels is not initialised (nil). ", ID)
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	sensChan, exists := atp.comms.sensorsChannels[ID]
	if !exists {
		atp.lock.comms.sensorsChannels.Unlock()
		err := errors.Errorf("Attempt to close sensor channel from real time (ID %d). ID doesn't exist", ID)
		atp.log.Warning(fmt.Sprintf("%+v", err))
		return err
	}

	close(sensChan.stop)
	atp.lock.comms.sensorsChannels.Unlock()
	atp.log.Info(fmt.Sprintf("Closing sensor channel (ID %d) signal sent", ID))
	return nil
}

// Start allows train movement. Once this method is called, Train Sensors are
// frequently updated according to the provided Setpoint.
// Returns an error if there are no active setpoint channels.
func (atp *Atp) Start() (<-chan struct{}, error) {
	ticker := time.NewTicker(time.Duration(200) * time.Millisecond)
	stop := make(chan struct{})

	go func() {
	loop:
		for {
			select {
			case <-ticker.C:
				// Check setpoint.Time is not too old
				a, _ := atp.core.GetSensors()
				atp.core.UpdateSensors(atp.getSetpoint(), time.Since(a.When()))
			case <-stop:
				ticker.Stop()
				break loop
			}
		}
	}()

	return stop, nil
}
