package downloadutils

import (
	"context"
	"github.com/bogem/id3v2/v2"
	"github.com/schollz/progressbar/v3"
	ffmpeggo "github.com/u2takey/ffmpeg-go"
	"github.com/wader/goutubedl"
	"github.com/zmb3/spotify/v2"
	"io"
	"log"
)

func DownloadSong(spotifyTrack *spotify.FullTrack, youtubeId string, outputPath string) {
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
	ffmpegCmd := ffmpeggo.Input("pipe:").Audio().Output(outputPath).Compile()

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

	// Set metadata
	tag, err := id3v2.Open(outputPath, id3v2.Options{Parse: false})

	artists := ""
	for i, artist := range spotifyTrack.Artists {
		if i == len(spotifyTrack.Artists)-1 {
			artists += artist.Name
			continue
		}

		artists += artist.Name + ", "
	}

	tag.SetArtist(artists)
	tag.SetTitle(spotifyTrack.Name)
	tag.SetAlbum(spotifyTrack.Album.Name)
	tag.SetYear(spotifyTrack.Album.ReleaseDate)
	tag.Save()
}
