package bparser_test

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/davidhintelmann/blockchain/bparser"
)

const (
	blocksFilePath = "C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data\\"
	blk00000Height = 119_965
)

func BenchmarkReadBlk(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode.")
	}

	matches, err := filepath.Glob(blocksFilePath + "*.dat")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	for i := 0; i < b.N; i++ {
		file, err := os.Open(matches[0])
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		defer file.Close()

		_, err = io.ReadAll(file)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}

func BenchmarkParseBlk(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode.")
	}

	matches, err := filepath.Glob(blocksFilePath + "*.dat")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	for i := 0; i < b.N; i++ {
		file, err := os.Open(matches[0])
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
		defer file.Close()

		readAll, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}

		// blkFileLen := len(readAll)
		// p := message.NewPrinter(language.English)
		// p.Printf("block file length: %d\n", blkFileLen)
		// p.Printf("block height: %d\n", blk00000Height)
		// fmt.Println()

		_, err = bparser.ParseBlocks(readAll, 0, 1, []byte{0})
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}
}
