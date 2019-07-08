package interfaces

import "time"

// Train contains data about train specifications.
type Train interface{}

// Track contains data about track specifications.
type Track interface{}

// Sensors contains dynamic data collected by train's sensors.
type Sensors interface {
	// When returns the time when the sensors were recorded.
	When() time.Time
}

// InitConditions contains initial conditions data of a Train.
type InitConditions interface{}

// Setpoint contains setpoint data.
type Setpoint interface {
	// When returns the time when the setpoint was recorded.
	When() time.Time
}

// Core collects essential information about train automation and
// implements train movement.
type Core interface {
	// GetTrain returns current Train specifications.
	GetTrain() (Train, error)

	// SetTrain sets new Train specifications.
	SetTrain(train Train) error

	// SetTrainInitConditions sets initial conditions of the Train
	// (i.e. position, velocity, acceleration, etc.)
	SetInitConditions(conditions InitConditions) error

	// GetTrack returns Track specifications by its ID.
	GetTrack(position int) (Track, error)

	// SetTracks sets an ordered slice of Track to drive through.
	SetTracks(tracks ...Track) error

	// GetSensors returns current sensors data.
	GetSensors() (Sensors, error)

	// SetSensors sets new sensor values.
	SetSensors(sensors Sensors) error

	// UpdateSensors updates Sensors from a setpoint and an elapsed time
	UpdateSensors(setpoint Setpoint, elapsed time.Duration) (Sensors, error)
}

// ATP provides a security layer over the Train movement implemented in Core.
// It enables setting setpoints and reading Sensor measurements.
type ATP interface {
	// GetTrain returns Train specifications.
	GetTrain() (Train, error)

	// SetTrain sets new Train specifications.
	SetTrain(train Train) error

	// GetTrack returns Track specifications by its ID.
	GetTrack(position int) (Track, error)

	// SetTracks sets an ordered slice of Track to drive through.
	SetTracks(track ...Track) error

	// SetInitConditions sets initial conditions of the Train
	// (i.e. position, velocity, acceleration, etc.)
	SetInitConditions(conditions InitConditions) error

	// OpenSetpointChannel creates a channel to deliver setpoints while driving.
	OpenSetpointChannel() (chan<- Setpoint, <-chan struct{}, error)

	// StopSetpointChannel stops ATP from reading setpoints and frees resources.
	StopSetpointChannel() error

	// NewSensorChannel creates a channel that delivers RealTime sensors data
	// at the given frequency rate.
	NewSensorChannel(ID int, frequency int) (<-chan Sensors, error)

	// CloseSensorChannel closes the SensorsChannel matching the given ID.
	CloseSensorChannel(ID int) error

	// Start makes the train move.
	Start() (<-chan struct{}, error)
}
