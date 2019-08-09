package core

import "github.com/trencat/Trencat/train/interfaces"

// ToInterfaceTracks converts []Core.Track to []interfaces.Track.
func ToInterfaceTracks(tracks ...Track) []interfaces.Track {
	interfacesTracks := make([]interfaces.Track, len(tracks))

	for i, track := range tracks {
		interfacesTracks[i] = track
	}

	return interfacesTracks
}
