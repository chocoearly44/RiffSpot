package spotify

import (
	"context"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"log"
)

func GetTracks(playlistId string) []spotify.PlaylistItemTrack {
	var toReturn []spotify.PlaylistItemTrack

	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     "0de60a2720584b1bb3d83b4ed48fa395",
		ClientSecret: "434adce2a93f4d5197f4f2565d0fcdfd",
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	tracks, err := client.GetPlaylistItems(
		ctx,
		spotify.ID(playlistId),
	)
	if err != nil {
		log.Fatal(err)
	}

	for page := 1; ; page++ {
		for _, item := range tracks.Items {
			if item.Track.Track == nil {
				continue
			}

			toReturn = append(toReturn, item.Track)
		}

		err = client.NextPage(ctx, tracks)
		if err == spotify.ErrNoMorePages {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	return toReturn
}
