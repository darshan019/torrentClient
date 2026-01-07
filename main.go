package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strings"

	// "log"
	"os"
	"torrentClient/file"
)

func shortHex(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func formatTrackerURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	q := u.Query()

	infoHash := q.Get("info_hash")
	peerID := q.Get("peer_id")

	return strings.Join([]string{
		"Tracker Endpoint : " + u.Scheme + "://" + u.Host + u.Path,
		"Info Hash        : " + shortHex(hex.EncodeToString([]byte(infoHash)), 16),
		"Peer ID          : " + shortHex(hex.EncodeToString([]byte(peerID)), 16),
		"Port             : " + q.Get("port"),
		"Left             : " + q.Get("left"),
		"Uploaded         : " + q.Get("uploaded"),
		"Downloaded       : " + q.Get("downloaded"),
		"Compact          : " + q.Get("compact"),
	}, "\n")
}

func humanBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func shortPeerID(id [20]byte) string {
	return fmt.Sprintf("%x", id[:6])
}

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

	// fmt.Println("Tracker URL: ")
	// fmt.Println(trackerURL)

	outPath := os.Args[1]

	/* ===============================
	   TORRENT INFO
	================================ */

	fmt.Println("=========== TORRENT INFO ===========")
	fmt.Printf(" Name           : %s\n", torrent.Name)
	fmt.Printf(" Size           : %s\n", humanBytes(torrent.Length))
	fmt.Printf(" Piece Length   : %s\n", humanBytes(torrent.PieceLength))
	fmt.Printf(" Pieces         : %d\n", len(torrent.PieceHashes))
	fmt.Printf(" Tracker        : %s\n", torrent.Announce)
	fmt.Printf(" Peer ID        : %s...\n", shortPeerID(peerID))
	fmt.Printf(" Output File    : %s\n", outPath)

	fmt.Println("\n=========== TRACKER REQUEST ==========")
	fmt.Println(formatTrackerURL(trackerURL))
	fmt.Println("=====================================")

	fmt.Println("================================")
	fmt.Println(" Starting download...")

	/* ===============================
	   DOWNLOAD
	================================ */

	err = torrent.DownloadToFile(outPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n================================")
	fmt.Println(" Download completed successfully")
	fmt.Println("================================")

}
