package matchingutils

import (
	"fmt"
	"github.com/raitonoberu/ytmusic"
	"github.com/zmb3/spotify/v2"
	"log"
	"math"
	"strings"
)

func FindClosestMatch(spotifyTrack *spotify.PlaylistItemTrack) *ytmusic.TrackItem {
	searchString := fmt.Sprintf(
		"%s - %s",
		spotifyTrack.Track.Name,
		spotifyTrack.Track.Artists[0].Name,
	)

	search := ytmusic.Search(searchString)
	matches, err := search.Next()
	if err != nil {
		log.Fatal(err)
	}

	tracks := make(map[*ytmusic.TrackItem]float64)
	var closestTrack *ytmusic.TrackItem

	for _, result := range matches.Tracks {
		tracks[result] = 0

		// Check for common word
		resultName := strings.Split(
			strings.Replace(strings.ToLower(result.Title), "[-/]", "", -1),
			" ",
		)

		trackNameWords := strings.Split(
			strings.Replace(strings.ToLower(spotifyTrack.Track.Name), "[-/]", "", -1),
			" ",
		)

		if len(resultName) != len(trackNameWords) {
			continue
		}

		commonWord := false
		for i, str := range resultName {
			if str == trackNameWords[i] {
				commonWord = true
				break
			}
		}

		if !commonWord {
			continue
		}

		// Check duration
		difference := math.Abs(float64(result.Duration - spotifyTrack.Track.Duration))
		nonMatchValue := (difference * difference) / float64(spotifyTrack.Track.Duration)
		durationMatch := 100 - (nonMatchValue * 100)

		tracks[result] += durationMatch
		closestTrack = result
	}

	for track, score := range tracks {
		if score > tracks[track] {
			closestTrack = track
		}
	}

	return closestTrack
}
