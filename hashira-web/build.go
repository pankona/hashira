package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	var (
		flagWatch bool
		flagClean bool
	)

	flag.BoolVar(&flagWatch, "watch", false, "watch the directory and build on change")
	flag.BoolVar(&flagClean, "clean", false, "clean artifacts")

	flag.Parse()

	outdir := filepath.Join("public", "assets")

	if flagClean {
		log.Printf("clean started\n")
		if err := cleanArtifacts(outdir); err != nil {
			log.Printf("clean failed: %v", err)
			os.Exit(1)
		}
		log.Printf("clean completed\n")
		return
	}

	buildOptions := api.BuildOptions{
		EntryPoints: []string{
			filepath.Join("src", "index.tsx"),
		},
		Target:   api.ES2015,
		Bundle:   true,
		Write:    true,
		Platform: api.PlatformBrowser,
		Outfile:  filepath.Join(outdir, "bundle.js"),
		Engines: []api.Engine{
			{Name: api.EngineChrome, Version: "58"},
			{Name: api.EngineFirefox, Version: "57"},
		},
	}

	if flagWatch {
		watch(buildOptions)
		return
	}

	build(buildOptions)
}

func cleanArtifacts(dir string) error {
	return os.RemoveAll(dir)
}

func watch(opts api.BuildOptions) error {
	opts.Watch = &api.WatchMode{
		OnRebuild: func(result api.BuildResult) {
			if len(result.Errors) > 0 {
				log.Printf("build error: %+v", result.Errors)
				return
			}
			log.Printf("rebuild succeeded")
		},
	}

	if err := build(opts); err != nil {
		log.Printf("failed to build: %v", err)
	}

	log.Println("start to watch the directory")

	// keep watching
	done := make(chan struct{})
	<-done

	return nil
}

func build(opts api.BuildOptions) error {
	log.Printf("build started\n")

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		log.Printf("%+v\n", result.Errors)
		os.Exit(1)
	}
	log.Printf("build completed\n")

	return nil
}
