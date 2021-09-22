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
		serve bool
		clean bool
	)

	flag.BoolVar(&serve, "serve", false, "start development server")
	flag.BoolVar(&clean, "clean", false, "clean artifacts")

	flag.Parse()

	outdir := filepath.Join("public", "assets")

	if clean {
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
			{Name: api.EngineSafari, Version: "11"},
			{Name: api.EngineEdge, Version: "16"},
		},
	}

	if serve {
		startDevServer(buildOptions)
		return
	}

	build(buildOptions)

	done := make(chan struct{})
	<-done
}

func cleanArtifacts(dir string) error {
	return os.RemoveAll(dir)
}

func startDevServer(opts api.BuildOptions) error {
	log.Println("starting dev server")

	server, err := api.Serve(api.ServeOptions{
		Servedir: filepath.Join("public"),
		Port:     8000,
	}, opts)
	if err != nil {
		log.Printf("failed to start development server: %v", err)
		os.Exit(1)
	}

	log.Printf("development server started with: %v:%v", server.Host, server.Port)

	return nil
}

func build(opts api.BuildOptions) error {
	log.Printf("build started\n")

	opts.Watch = &api.WatchMode{
		OnRebuild: func(result api.BuildResult) {
			if len(result.Errors) > 0 {
				log.Printf("build error: %+v", result.Errors)
				return
			}
		},
	}

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		log.Printf("%+v\n", result.Errors)
		os.Exit(1)
	}
	log.Printf("build completed\n")

	return nil
}
