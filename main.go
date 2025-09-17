package main

import (
	"fileman/clock"
	"fileman/config"
	"fileman/fs"
	"fileman/handler"
	"github.com/go-co-op/gocron/v2"
	"os"
	"runtime"
)

func main() {
	realClock := clock.RealClock{}
	fileHandler := handler.New(realClock)
	fileSystem := fs.FS{}

	logger := gocron.NewLogger(gocron.LogLevelInfo)
	scheduler, err := gocron.NewScheduler(gocron.WithLogger(logger))
	if err != nil {
		panic(err)
	}

	configFile, configExists := os.LookupEnv("CONFIG_PATH")
	if !configExists {
		configFile = "config.json"
	}

	configObject, configError := config.New(configFile).Load()
	if configError != nil {
		panic(configError)
	}

	if configObject.Cron == "" {
		panic("Cron expression not set in configObject")
	}

	errs := make([]error, 0)
	jobs := make([]gocron.Job, 0)

	for _, directory := range configObject.WatchedDirectories {
		job, e := scheduler.NewJob(
			gocron.CronJob(configObject.Cron, false),
			gocron.NewTask(func() {
				deleted, errs := fileHandler.DeleteOldFiles(fileSystem, directory.Path, directory.Age)
				for _, d := range deleted {
					logger.Info("Deleted file", d)
				}

				for _, e := range errs {
					logger.Error("Error deleting file", e.Error())
				}

				if len(deleted) == 0 && len(errs) == 0 {
					logger.Info("No files to delete in path", directory.Path)
				}
			}),
			gocron.WithName("PathCleaner-"+directory.Path),
		)

		if e != nil {
			errs = append(errs, e)
			continue
		}

		jobs = append(jobs, job)
	}

	for _, e := range errs {
		logger.Error("Error scheduling job", e.Error())
	}

	for _, job := range jobs {
		logger.Info("Scheduled Job", "Name", job.Name(), "ID", job.ID())
	}

	scheduler.Start()

	runtime.Goexit()
}
