package downloadutils

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/raitonoberu/ytmusic"
	"github.com/schollz/progressbar/v3"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/wader/goutubedl"
	"github.com/zmb3/spotify/v2"
)

func FindClosestMatch(spotifyTrack *spotify.PlaylistItemTrack) *ytmusic.TrackItem {
	searchString := fmt.Sprintf(
		"%s - %s",
		spotifyTrack.Track.Name,
		spotifyTrack.Track.Artists[0].Name,
	)

	search := ytmusic.Search(searchString)
	result, err := search.Next()
	if err != nil {
		log.Fatal(err)
	}

	minDiff := math.Abs(float64(spotifyTrack.Track.Duration - result.Tracks[0].Duration))
	closestTrack := result.Tracks[0]

	for _, track := range result.Tracks {
		diff := math.Abs(float64(spotifyTrack.Track.Duration - track.Duration))

		if diff < minDiff {
			minDiff = diff
			closestTrack = track
		}
	}

	return closestTrack
}

func DownloadSong(youtubeId string, outputPath string) {
	log.SetOutput(io.Discard)
	goutubedl.Path = "yt-dlp"

	// Create new youtube download
	result, err := goutubedl.New(
		context.Background(),
		youtubeId,
		goutubedl.Options{Type: goutubedl.TypeSingle},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start downloading the video
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		log.Fatal(err)
	}

	// Create a progress bar
	progressBar := progressbar.DefaultBytes(
		int64(result.Info.FilesizeApprox),
		"Downloading "+result.Info.Title,
	)
	defer progressBar.Close()

	// Create ffmpeg audio extractor
	ffmpegCmd := ffmpeg_go.Input("pipe:").Audio().Output(outputPath).Compile()

	// Retrieve the standard input pipe of ffmpeg
	ffmpegStdin, err := ffmpegCmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Run the ffmpeg command in background
	err = ffmpegCmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Pipe youtube download stream to progress bar and ffmpeg
	io.Copy(io.MultiWriter(progressBar, ffmpegStdin), downloadResult)

}
