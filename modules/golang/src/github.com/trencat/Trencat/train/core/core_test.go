package core_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
	"testing"
	"time"

	"github.com/trencat/Trencat/train/core"
	"github.com/trencat/Trencat/train/interfaces"
)

var flagUpdate bool
var co interfaces.Core
var log *syslog.Writer

type testdataTrain map[string]core.Train
type testdataTrack map[string][]core.Track
type testdataScenario struct {
	TestTrain string
	TestTrack string
	Sensors   core.Sensors
	Setpoint  core.Setpoint
	Duration  time.Duration
	Expected  core.Sensors
}
type testdataUpdateSensors map[string]testdataScenario

// unmarshalFromFile decodes a json file into v variable.
func unmarshalFromFile(path string, v interface{}, t *testing.T) {
	t.Helper()

	// read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// unmarshal
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("%+v", err)
	}
}

// unarmshalFromFileKey decodes a json file, gets the specified key and
// stores it into v variable.
func unmarshalFromFileKey(path string, key string, v interface{}, t *testing.T) {
	t.Helper()

	switch x := v.(type) {
	case *core.Train:
		testdata := make(testdataTrain)
		unmarshalFromFile(path, &testdata, t)
		value, exists := testdata[key]
		if !exists {
			t.Fatalf("Key %s does not exist in %s", key, path)
		}
		*x = value
	case *[]core.Track:
		testdata := make(testdataTrack)
		unmarshalFromFile(path, &testdata, t)
		value, exists := testdata[key]
		if !exists {
			t.Fatalf("Key %s does not exist in %s", key, path)
		}
		*x = value
	}
}

// marshalToFile encodes a mapping into a json file.
func marshalToFile(path string, v interface{}, t *testing.T) {
	t.Helper()

	// marshal
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// write file
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

// setScenario sets a train, tracks and initial conditions to a core.Core instance.
func setScenario(co *core.Core, scenario *testdataScenario, t *testing.T) {
	t.Helper()

	var train core.Train
	var tracks []core.Track

	unmarshalFromFileKey("testdata/trains.json", scenario.TestTrain, &train, t)
	unmarshalFromFileKey("testdata/tracks.json", scenario.TestTrack, &tracks, t)

	err := co.SetTrain(train)
	if err != nil {
		t.Fatalf("Cannot set Train%+v. Got error: %s", train, err)
	}

	err = co.SetTracks(core.ToInterfaceTracks(tracks...)...)
	if err != nil {
		t.Fatalf("Cannot set Track%+v. Got error: %s", tracks, err)
	}

	err = co.SetInitConditions(scenario.Sensors)
	if err != nil {
		t.Fatalf("Cannot set init conditions Sensors%+v", scenario.Sensors)
	}
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
	testdata := make(testdataTrain)
	unmarshalFromFile("testdata/trains.json", &testdata, t)

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
	testdata := make(testdataTrack)
	unmarshalFromFile("testdata/tracks.json", &testdata, t)

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
	unmarshalFromFile("testdata/updateSensorsAcceleration.json", &testdata, t)

	for alias, scenario := range testdata {

		setScenario(co.(*core.Core), &scenario, t)

		newSensor, err := co.UpdateSensors(scenario.Setpoint, scenario.Duration)
		if err != nil {
			t.Errorf("With scenario %s,\nGot error %s.\nExpected nil", alias, err)
			continue
		}

		if flagUpdate {
			// Update scenario expected value
			scenario.Expected = newSensor.(core.Sensors)
			// Update testdata value
			testdataScenario := testdata[alias]
			testdataScenario.Expected = newSensor.(core.Sensors)
			testdata[alias] = testdataScenario
		}

		if scenario.Expected != newSensor {
			t.Errorf("With scenario %s,\nGot Sensors%+v,\nExpected Sensors%+v", alias, scenario.Expected, newSensor)
			continue
		}

	}

	if flagUpdate {
		marshalToFile("testdata/updateSensorsAcceleration.json", testdata, t)
	}
}
