package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"torrentClient/file"
)

func main() {
	fmt.Println("hello")

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

	fmt.Println("Announce URL:", torrent.Announce)
	fmt.Println("File Name:", torrent.Name)
	fmt.Println("File Length:", torrent.Length)
	fmt.Println("Piece Length:", torrent.PieceLength)
	fmt.Println("Info Hash:", torrent.InfoHash)
	fmt.Println("Number of Pieces:", len(torrent.PieceHashes))

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

	peersFound, err := torrent.RequestPeers(peerID, 6881)
	if err != nil {
		panic(err)
	}
	fmt.Println("Peers:", peersFound)
}
