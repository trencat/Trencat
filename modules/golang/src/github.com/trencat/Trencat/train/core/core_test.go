package core_test

import (
	"flag"
	"fmt"
	"log/syslog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/trencat/Trencat/testutils"
	"github.com/trencat/Trencat/train/core"
	"github.com/trencat/Trencat/train/interfaces"
)

// TODO: Read these constants from file?
var testdataTrainsPath string = filepath.Join("..", "..", "testutils", "testdata", "trains.json")
var testdataTracksPath = filepath.Join("..", "..", "testutils", "testdata", "tracks.json")
var testdataScenariosPath = filepath.Join("..", "..", "testutils", "testdata", "scenarios.json")
var testdataPath = filepath.Join("testdata", "updateSensorsAcceleration.json")

var flagUpdate bool
var co interfaces.Core
var log *syslog.Writer

type testdataUpdateSensors map[string]struct {
	Scenario string
	Setpoint core.Setpoint
	Duration time.Duration
	Expected core.Sensors
}

func TestMain(m *testing.M) {
	// Parse arguments
	flag.BoolVar(&flagUpdate, "update", false, "Update golden file tests.")
	flag.Parse()

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
	testdata := make(testutils.TestdataTrain)
	testutils.UnmarshalFromFile(testdataTrainsPath, &testdata, t)

	for alias, train := range testdata {
		error := co.SetTrain(train)
		if error != nil {
			t.Errorf("With %s Train%+v, got error %s. Expected nil", alias, train, error)
			continue
		}
		coreTrain, error := co.GetTrain()
		if error != nil {
			t.Errorf("With %s Train%+v, got error %s. Expected nil", alias, train, error)
			continue
		}
		if coreTrain != train {
			t.Errorf("With %s got Train%+v, expected Train%+v", alias, coreTrain, train)
			continue
		}
	}
}

func TestSetGetTrack(t *testing.T) {
	testdata := make(testutils.TestdataTrack)
	testutils.UnmarshalFromFile(testdataTracksPath, &testdata, t)

	for alias, tracks := range testdata {
		// Convert []core.Track to []interfaces.Track
		interfaceTracks := core.ToInterfaceTracks(tracks...)

		error := co.SetTracks(interfaceTracks...)
		if error != nil {
			t.Errorf("With %s []Track %+v, got error %s. Expected nil", alias, tracks, error)
			continue
		}

		// Check all tracks
		for i := 0; i < len(tracks); i++ {
			track := tracks[i]
			coreTrack, error := co.GetTrack(i)
			if error != nil {
				t.Errorf("With %s[%d] Track%+v, got error %s. Expected nil", alias, i, track, error)
				continue
			}

			if coreTrack != track {
				t.Errorf("With %s[%d] got Track%+v, expected Track%+v", alias, i, coreTrack, track)
				continue
			}
		}
	}
}

func TestSetGetInitConditions(t *testing.T) {
	//TODO
}

// TestUpdateSensorsAcceleration tests UpdateSensors implementation
// considering that setpoint refers to acceleration.
func TestUpdateSensorsAcceleration(t *testing.T) {
	testdata := make(testdataUpdateSensors)
	testutils.UnmarshalFromFile(testdataPath, &testdata, t)

	for alias, test := range testdata {
		//Read scenario
		scenario := testutils.Scenario{}
		testutils.UnmarshalFromFileKey(testdataScenariosPath, test.Scenario, &scenario, t)

		testutils.SetTrainScenario(filepath.Dir(testdataScenariosPath), co.(*core.Core), &scenario, t)

		newSensor, err := co.UpdateSensors(test.Setpoint, test.Duration)
		if err != nil {
			t.Errorf("With scenario %s,\nGot error %s.\nExpected nil", alias, err)
			continue
		}

		if flagUpdate {
			// Update scenario expected value
			test.Expected = newSensor.(core.Sensors)
			// Update testdata value
			testdataScenario := testdata[alias]
			testdataScenario.Expected = newSensor.(core.Sensors)
			testdata[alias] = testdataScenario
		}

		if test.Expected != newSensor {
			t.Errorf("With scenario %s,\nGot Sensors%+v,\nExpected Sensors%+v", alias, test.Expected, newSensor)
			continue
		}
	}

	if flagUpdate {
		testutils.MarshalToFile("testdata/updateSensorsAcceleration.json", testdata, t)
	}
}
