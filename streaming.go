package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

const PLAYLIST_FOLDER = "./content/stream/"
const SEGMENTS_FOLDER = PLAYLIST_FOLDER + "segments/"

const PLAYLIST_FILE = "segment.m3u8"
const STREAM_FILE_TEMPLATE = "segment?.ts"
const DEFAULT_SEGMENT_COUNT = 3

const STATIC_HEADER = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-MEDIA-SEQUENCE:`

const SEGMENT_TEMPLATE = `#EXTINF:10.000000
segments/segment%d.ts
`

const STATIC_ENDING = "#EXT-X-ENDLIST"

var counter = 0
var maxCounter int

func initMaxCounter() {
	if total, err := getTotalSegments(SEGMENTS_FOLDER); err != nil {
		log.Fatal(errors.Wrap(err, "Error initiating maxCounter"))
	} else {
		maxCounter = total - DEFAULT_SEGMENT_COUNT
	}
}

func getTotalSegments(folder string) (int, error) {
	if files, err := os.ReadDir(folder); err != nil {
		return 0, errors.Wrap(err, "Error reading total segments")
	} else {
		return len(files), nil

	}
}
