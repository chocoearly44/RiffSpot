package playlist

import (
	"RiffSpot/internal/providers/spotify"
	"RiffSpot/internal/utils/downloadutils"
	"RiffSpot/internal/utils/matchingutils"
	"fmt"
	"path/filepath"
)

func Execute(playlistId string, outputFolder string) {
	for _, spotifyTrack := range spotify.GetTracks(playlistId) {
		video := matchingutils.FindClosestMatch(&spotifyTrack)
		if video == nil {
			fmt.Printf("\n%s not found.\n", spotifyTrack.Track.Name)
			continue
		}
		go downloadutils.DownloadSong(spotifyTrack.Track, video.VideoID, filepath.Join(outputFolder, spotifyTrack.Track.Name+".mp3"))
	}
}
