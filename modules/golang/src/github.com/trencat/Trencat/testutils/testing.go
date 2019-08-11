package testutils

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/trencat/Trencat/train/core"
)

type TestdataTrain map[string]core.Train
type TestdataTrack map[string][]core.Track
type Scenario struct {
	TestTrain string
	TestTrack string
	Sensors   core.Sensors
}
type TestdataScenario map[string]Scenario

// UnmarshalFromFile decodes a json file into v variable.
func UnmarshalFromFile(path string, v interface{}, t *testing.T) {
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

// UnmarshalFromFileKey decodes a json file, gets the specified key and
// stores it into v variable.
func UnmarshalFromFileKey(path string, key string, v interface{}, t *testing.T) {
	t.Helper()

	switch x := v.(type) {
	case *core.Train:
		testdata := make(TestdataTrain)
		UnmarshalFromFile(path, &testdata, t)
		value, exists := testdata[key]
		if !exists {
			t.Fatalf("Key %s does not exist in %s", key, path)
		}
		*x = value
	case *[]core.Track:
		testdata := make(TestdataTrack)
		UnmarshalFromFile(path, &testdata, t)
		value, exists := testdata[key]
		if !exists {
			t.Fatalf("Key %s does not exist in %s", key, path)
		}
		*x = value
	case *Scenario:
		testdata := make(TestdataScenario)
		UnmarshalFromFile(path, &testdata, t)
		value, exists := testdata[key]
		if !exists {
			t.Fatalf("Key %s does not exist in %s", key, path)
		}
		*x = value
	default:
		t.Fatalf("Cannot unmarshal %s", path)
	}
}

// MarshalToFile encodes a mapping into a json file.
func MarshalToFile(path string, v interface{}, t *testing.T) {
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

// SetScenario sets a train, tracks and initial conditions to a core.Core instance.
func SetScenario(co *core.Core, scenario *Scenario, t *testing.T) {
	t.Helper()

	var train core.Train
	var tracks []core.Track

	UnmarshalFromFileKey("testdata/trains.json", scenario.TestTrain, &train, t)
	UnmarshalFromFileKey("testdata/tracks.json", scenario.TestTrack, &tracks, t)

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
