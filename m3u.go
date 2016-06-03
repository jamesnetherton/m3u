package m3u

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Playlist is a type that represents an m3u playlist containing 0 or more tracks
type Playlist struct {
	Tracks []Track
}

// Track represents an m3u track
type Track struct {
	Name   string
	Length int
	Path   string
}

// Parse parses an m3u playlist with the given file name and returns a Playlist
func Parse(fileName string) (playlist Playlist, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		err = errors.New("Unable to open playlist file")
		return
	}

	lineCount := 1
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if lineCount == 1 && !strings.HasPrefix(line, "#EXTM3U") {
			err = errors.New("Invalid m3u file format. Expected #EXTM3U file header")
			return
		}

		lineCount++

		if strings.HasPrefix(line, "#EXTINF") {
			line := strings.Replace(line, "#EXTINF:", "", -1)
			trackInfo := strings.Split(line, ",")
			if len(trackInfo) < 2 {
				err = errors.New("Invalid m3u file format. Expected EXTINF metadata to contain track length and name data")
				return
			}
			length, _ := strconv.Atoi(trackInfo[0])
			track := &Track{trackInfo[1], length, ""}
			playlist.Tracks = append(playlist.Tracks, *track)
		} else if strings.HasPrefix(line, "#") {
			continue
		} else {
			playlist.Tracks[len(playlist.Tracks)-1].Path = line
		}
	}

	return playlist, nil
}
