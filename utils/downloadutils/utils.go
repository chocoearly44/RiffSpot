package downloadutils

import (
	"context"
	"flag"
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/wader/goutubedl"
	"log"
)

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
