package riffspot

import (
	"RiffSpot/internal/commands/playlist"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var outputFolder string

var rootCmd = &cobra.Command{
	Use:   "riffspot",
	Short: "RiffSpot - Blazingly fast Spotify downloader.",
	Long:  "RiffSpot is the best Spotify playlist downloader on this side of the Mississippi.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "Download playlists",
	Long:  "Use this command when downloading playlists.",
	Run: func(cmd *cobra.Command, args []string) {
		playlist.Execute(args[0], cmd.Flag("output").Value.String())
	},
}

func init() {
	playlistCmd.Flags().StringVarP(&outputFolder, "output", "o", "", "Save songs to folder")
	playlistCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(playlistCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
