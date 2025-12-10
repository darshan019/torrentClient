package main

import (
	"crypto/rand"
	"fmt"
	"log"

	// "log"
	"os"
	"torrentClient/file"
)

func main() {
	f, err := os.Open("debian.torrent")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bto, err := file.Open(f)
	if err != nil {
		panic(err)
	}

	torrent, err := bto.ToTorrentFile()
	if err != nil {
		panic(err)
	}

	var peerID [20]byte
	_, err = rand.Read(peerID[:])
	if err != nil {
		panic(err)
	}

	trackerURL, err := torrent.BuildTrackerURL(peerID, 6881)
	if err != nil {
		panic(err)
	}

	fmt.Println("Tracker URL: ")
	fmt.Println(trackerURL)

	outPath := os.Args[1]
	err = torrent.DownloadToFile(outPath)
	if err != nil {
		log.Fatal(err)
	}

}
