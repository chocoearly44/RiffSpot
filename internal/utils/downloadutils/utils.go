package downloadutils

import (
	"context"
	"flag"
	"fmt"
	"github.com/raitonoberu/ytmusic"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/wader/goutubedl"
	"github.com/zmb3/spotify/v2"
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
	goutubedl.Path = "yt-dlp"

	log.SetFlags(0)
	flag.Parse()

	result, err := goutubedl.New(
		context.Background(),
		youtubeId,
		goutubedl.Options{Type: goutubedl.TypeSingle, DebugLog: log.Default()},
	)
	if err != nil {
		log.Fatal(err)
	}

	filter := "best"

	dr, err := result.Download(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	com := ffmpeg_go.Input("pipe:").Audio().Output(path).Compile()
	com.Stdin = dr

	err = com.Run()
	if err != nil {
		fmt.Println("Error transcoding")
	}
}
