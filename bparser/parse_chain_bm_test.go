package bparser_test

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	blocksFilePath = "C:\\Users\\david\\OneDrive\\Documents\\code\\python\\Blockchain\\Bitcoin\\data\\bitcoin_data\\"
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
