package main

import (
	"Galia/cluster"
	"Galia/conf"
	"Galia/console"
	"Galia/log"
	"Galia/module"
	"os"
	"os/signal"
)

func Run(mods ...module.Module) {
	// logger

	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath, "", conf.LogFlag, 10)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Galia %v starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
