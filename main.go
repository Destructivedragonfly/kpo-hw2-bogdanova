package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"example.com/pos/adapters"
	"example.com/pos/application"
	"example.com/pos/domain"
)

func main() {
	dataPath := flag.String("data", "pos-data", "directory for JSON data files")
	demo := flag.Bool("demo", false, "run a scripted end-to-end demonstration")
	flag.Parse()

	if *demo {
		dir, err := os.MkdirTemp("", "pos-demo-")
		if err != nil {
			fatal(err)
		}
		defer os.RemoveAll(dir)
		*dataPath = filepath.Join(dir, "data")
	}

	// DI происходит тут
	store, err := adapters.NewJSONFileStore(*dataPath)
	if err != nil {
		fatal(err)
	}
	ids := adapters.NewSequentialIDGenerator()
	finalGradeService, err := domain.NewFinalGradeService(domain.WeightedAverageStrategy{})
	if err != nil {
		fatal(err)
	}
	pos := application.NewPOSFacade(
		store, store, store, ids,
		domain.NewGradeFactory(), finalGradeService,
	)

	if *demo {
		if err := adapters.RunDemo(context.Background(), pos, os.Stdout); err != nil {
			fatal(err)
		}
		return
	}

	console := adapters.NewConsoleUI(pos)
	console.Run(context.Background())
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "POS:", err)
	os.Exit(1)
}
