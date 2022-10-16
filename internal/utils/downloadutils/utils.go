package downloadutils

import (
	"context"
	"fmt"
	"github.com/raitonoberu/ytmusic"
	"github.com/schollz/progressbar/v3"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/wader/goutubedl"
	"github.com/zmb3/spotify/v2"
	"io"
	"log"
	"math"
)

func FindClosestMatch(spotifyTrack *spotify.PlaylistItemTrack) *ytmusic.TrackItem {
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

	return closestTrack
}

func DownloadSong(youtubeId string, path string) {
	log.SetOutput(io.Discard)
	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(
		context.Background(),
		youtubeId,
		goutubedl.Options{Type: goutubedl.TypeSingle},
	)

	if err != nil {
		log.Fatal(err)
	}

	dr, err := result.Download(context.Background(), "best")
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.DefaultBytes(
		int64(result.Info.FilesizeApprox),
		"Downloading "+result.Info.Title,
	)

	com := ffmpeg_go.Input("pipe:").Audio().Output(path).Compile()
	com.Stdin = dr
	if err != nil {
		fmt.Println("Error transcoding")
	}
	io.Copy(bar, dr)

	err = com.Run()
	bar.Close()
}
