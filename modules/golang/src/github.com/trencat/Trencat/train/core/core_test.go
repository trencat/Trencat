package core_test

import (
	"fmt"
	"log/syslog"
	"os"
	"testing"

	"github.com/trencat/Trencat/train/core"
	"github.com/trencat/Trencat/train/interfaces"
)

var log *syslog.Writer
var co interfaces.Core

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	// Setup
	syslog, error := syslog.Dial("tcp", "localhost:514",
		syslog.LOG_WARNING|syslog.LOG_LOCAL0, "coreTest")

	if error != nil {
		panic(fmt.Sprintf("%s", error))
	}

	log = syslog

	testCore, err := core.New(log)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	co = &testCore

	//Teardown
	os.Exit(m.Run())
}

func TestSetGetTrain(t *testing.T) {
	f := core.NewFactory()
	for i := 0; i <= 10; i++ {
		train := f.GetTrain()
		error := co.SetTrain(train)
		if error != nil {
			t.Fatalf("With input Train%+v, got error %s. Expected nil", train, error)
		}
		coreTrain, error := co.GetTrain()
		if error != nil {
			t.Fatalf("With input Train%+v, got error %s. Expected nil", train, error)
		}
		if coreTrain != train {
			t.Fatalf("Got Train%+v, expected Train%+v", coreTrain, train)
		}
	}
}

func TestSetGetTrack(t *testing.T) {

	f := core.NewFactory()
	tracks, error := f.GetTrack(10, 500, 5000, true, true, true)
	if error != nil {
		t.Fatalf("%s", error)
	}

	// Test core.SetTracks
	co.SetTracks(tracks...)

	for i, track := range tracks {
		coreTrack, error := co.GetTrack(i)
		if error != nil {
			t.Fatalf("%s", error)
		}

		if coreTrack != track {
			t.Fatalf("Got Track%+v, expected Track%+v", coreTrack, track)
		}
	}
}

//TestUpdateSensorsAcceleration tests UpdateSensors implementation
//considering that setpoint refers to acceleration.
func TestUpdateSensorsAcceleration(t *testing.T) {
	//Create table driven tests
}
