package playlist

import (
	"RiffSpot/internal/providers/spotify"
	"RiffSpot/internal/utils/downloadutils"
	"path/filepath"
)

func Execute(playlistId string, outputFolder string) {
	for _, spotifyTrack := range spotify.GetTracks(playlistId) {
		videoId := downloadutils.FindClosestMatch(&spotifyTrack).VideoID
		downloadutils.DownloadSong(videoId, filepath.Join(outputFolder, spotifyTrack.Track.Name+".mp3"))
	}
}
