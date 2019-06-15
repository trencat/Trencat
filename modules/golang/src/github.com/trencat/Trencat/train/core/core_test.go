package core_test

import (
	"fmt"
	"log/syslog"
	"os"
	"testing"

	"github.com/trencat/Trencat/train/core"
)

var log *syslog.Writer

func TestMain(m *testing.M) {
	// Setup
	syslog, error := syslog.Dial("tcp", "localhost:514",
		syslog.LOG_WARNING|syslog.LOG_LOCAL0, "ATP_test")

	if error != nil {
		panic(fmt.Sprintf("%s", error))
	}

	log = syslog

	// call flag.Parse() here if TestMain uses flags

	//Teardown
	os.Exit(m.Run())
}

func TestSetGetTrain(t *testing.T) {
	co, error := core.New(log)
	if error != nil {
		t.Fatalf("%s", error)
	}

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
	co, error := core.New(log)
	if error != nil {
		t.Fatalf("%s", error)
	}

	f := core.NewFactory()
	tracks, error := f.GetTrack(5, 500, 5000, true, true, true)
	if error != nil {
		t.Fatalf("%s", error)
	}

	for _, track := range tracks {
		error := co.InsertTrack(track)
		if error != nil {
			t.Fatalf("%s", error)
		}
	}

	for _, track := range tracks {
		coreTrack, error := co.GetTrack(track.ID)
		if error != nil {
			t.Fatalf("%s", error)
		}

		if coreTrack != track {
			t.Fatalf("Got Track%+v, expected Track%+v", coreTrack, track)
		}
	}
}

func TestIncorrectTrack(t *testing.T) {
	co, error := core.New(log)
	if error != nil {
		t.Fatalf("%s", error)
	}

	f := core.NewFactory()
	tracks, error := f.GetTrack(2, 500, 5000, true, true, true)
	if error != nil {
		t.Fatalf("%s", error)
	}

	// Non matching tracks
	tracks[0].NextTrackID = tracks[1].ID - 1

	error = co.InsertTrack(tracks[0])
	if error != nil {
		t.Fatalf("%s", error)
	}

	//Should return an error
	error = co.InsertTrack(tracks[1])
	if error == nil {
		t.Fatal("Got nil error, epected non-nil error")
	}
}
