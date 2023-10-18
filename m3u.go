package m3u

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Playlist is a type that represents an m3u playlist containing 0 or more tracks or streams
type Playlist struct {
	Tracks  []Track
	Streams []StreamInfo
}

// A Tag is a simple key/value pair
type Tag struct {
	Name  string
	Value string
}

// Track represents an m3u track with a Name, Lengh, URI and a set of tags
type Track struct {
	Name   string
	Length int
	URI    string
	Tags   []Tag
}

type StreamInfo struct {
	Resolution string
	Bandwidth  int
	Codecs     string
	Name       string
	FrameRate  float64
	URI        string
}

// Parse parses an m3u playlist with the given file name and returns a Playlist
func Parse(fileName string) (Playlist, error) {
	var f io.ReadCloser

	if strings.HasPrefix(fileName, "http://") || strings.HasPrefix(fileName, "https://") {
		data, err := http.Get(fileName)
		if err != nil {
			return Playlist{},
				fmt.Errorf("unable to open playlist URL: %v", err)
		}
		f = data.Body
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			return Playlist{},
				fmt.Errorf("unable to open playlist file: %v", err)
		}
		f = file
	}
	defer f.Close()

	onFirstLine := true
	scanner := bufio.NewScanner(f)
	tagsRegExp, _ := regexp.Compile("([a-zA-Z0-9-]+?)=\"([^\"]+)\"")
	playlist := Playlist{}

	for scanner.Scan() {
		line := scanner.Text()
		if onFirstLine && !strings.HasPrefix(line, "#EXTM3U") {
			return Playlist{},
				errors.New("invalid m3u file format. Expected #EXTM3U file header")
		}

		onFirstLine = false

		if strings.HasPrefix(line, "#EXTINF") {
			line := strings.Replace(line, "#EXTINF:", "", -1)
			trackInfo := strings.Split(line, ",")
			if len(trackInfo) < 2 {
				return Playlist{},
					errors.New("invalid m3u file format. Expected EXTINF metadata to contain track length and name data")
			}
			length, parseErr := strconv.Atoi(strings.Split(trackInfo[0], " ")[0])
			if parseErr != nil {
				return Playlist{}, errors.New("unable to parse length")
			}
			track := &Track{strings.Trim(trackInfo[1], " "), length, "", nil}
			tagList := tagsRegExp.FindAllString(line, -1)
			for i := range tagList {
				tagInfo := strings.Split(tagList[i], "=")
				tag := &Tag{tagInfo[0], strings.Replace(tagInfo[1], "\"", "", -1)}
				track.Tags = append(track.Tags, *tag)
			}
			playlist.Tracks = append(playlist.Tracks, *track)
		} else if strings.HasPrefix(line, "#EXT-X-STREAM-INF") {
			line := strings.Replace(line, "#EXT-X-STREAM-INF:", "", -1)
			streamInfo := strings.Split(line, ",")
			if len(streamInfo) < 2 {
				return Playlist{},
					errors.New("invalid m3u file format. Expected EXT-X-STREAM-INF metadata to contain bitrate and resolution data")
			}
			bandwidth := strings.Split(streamInfo[1], "=")[1]
			bandwidthInt, parseErr := strconv.Atoi(bandwidth)
			if parseErr != nil {
				return Playlist{}, errors.New("unable to parse bandwidth")
			}
			resolution := strings.Split(streamInfo[2], "=")[1]
			codecs := strings.Split(streamInfo[3], "=")[1]
			frameRate := strings.Split(streamInfo[4], "=")[1]
			frameRateFloat, parseErr := strconv.ParseFloat(frameRate, 64)
			if parseErr != nil {
				return Playlist{}, errors.New("unable to parse frame rate")
			}
			stream := &StreamInfo{resolution, bandwidthInt, codecs, "", frameRateFloat, ""}
			playlist.Streams = append(playlist.Streams, *stream)
		} else if strings.HasPrefix(line, "#") || line == "" {
			continue
		} else if len(playlist.Tracks) == 0 && len(playlist.Streams) == 0 {
			return Playlist{},
				errors.New("URI provided for playlist with no tracks or streams")

		} else {
			playlist.Tracks[len(playlist.Tracks)-1].URI = strings.Trim(line, " ")
		}
	}

	return playlist, nil
}

// Marshall Playlist to an m3u file.
func Marshall(p Playlist) (io.Reader, error) {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	if err := MarshallInto(p, w); err != nil {
		return nil, err
	}

	return buf, nil
}

// MarshallInto a *bufio.Writer a Playlist.
func MarshallInto(p Playlist, into *bufio.Writer) error {
	into.WriteString("#EXTM3U\n")
	for _, track := range p.Tracks {
		into.WriteString("#EXTINF:")
		into.WriteString(fmt.Sprintf("%d ", track.Length))
		for i := range track.Tags {
			if i == len(track.Tags)-1 {
				into.WriteString(fmt.Sprintf("%s=%q", track.Tags[i].Name, track.Tags[i].Value))
				continue
			}
			into.WriteString(fmt.Sprintf("%s=%q ", track.Tags[i].Name, track.Tags[i].Value))
		}
		into.WriteString(", ")

		into.WriteString(fmt.Sprintf("%s\n%s\n", track.Name, track.URI))
	}

	return into.Flush()
}
