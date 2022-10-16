package main

import (
	"RiffSpot/providers/spotify"
	"RiffSpot/utils/downloadutils"
	"github.com/raitonoberu/ytmusic"
	"math"
)

func main() {
	for _, spotifyTrack := range spotify.GetTracks("75LakzUrwvcTaOQJFK70sj") {
		s := ytmusic.Search(spotifyTrack.Track.Name)
		result, err := s.Next()

		if err != nil {
			panic(err)
		}

		minDiff := math.Abs(float64(spotifyTrack.Track.Duration - result.Tracks[0].Duration))
		closestTrack := result.Tracks[0]

		for _, i := range result.Tracks {
			diff := math.Abs(float64(spotifyTrack.Track.Duration - i.Duration))

			if diff < minDiff {
				minDiff = diff
				closestTrack = i
			}
		}

		downloadutils.DownloadSong(closestTrack.VideoID, "/mnt/Data/Temp/Crinhen/"+spotifyTrack.Track.Name+".mp3")
	}
}
