package testutils

import (
	"path"
	"testing"

	"github.com/trencat/Trencat/train/core"
)

// SetTrainScenario sets a train, tracks and initial conditions to a core.Core instance.
func SetTrainScenario(testdataDir string, co *core.Core, scenario *Scenario, t *testing.T) {
	t.Helper()

	var train core.Train
	var tracks []core.Track

	UnmarshalFromFileKey(path.Join(testdataDir, "trains.json"), scenario.TestTrain, &train, t)
	UnmarshalFromFileKey(path.Join(testdataDir, "tracks.json"), scenario.TestTrack, &tracks, t)

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
