package main

import (
	"fileman/clock"
	"fileman/fs"
	"fileman/handler"
	"github.com/go-co-op/gocron/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
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

	crontab, crontabExists := os.LookupEnv("FILEMAN_CRONTAB")
	if !crontabExists {
		panic("FILEMAN_CRONTAB not set")
	}

	paths, pathsExists := os.LookupEnv("FILEMAN_PATHS")
	pathsSlice := strings.Split(paths, ",")
	if !pathsExists {
		panic("FILEMAN_PATHS not set")
	}

	ages, agesExists := os.LookupEnv("FILEMAN_AGES")
	if !agesExists {
		panic("FILEMAN_AGES not set")
	}

	agesSlice := strings.Split(ages, ",")
	if len(agesSlice) != len(pathsSlice) {
		panic("Number of ages doesn't match number of paths")
	}

	errs := make([]error, 0)
	jobs := make([]gocron.Job, 0)

	for idx, path := range pathsSlice {
		// Convert age from string to float
		age, e := strconv.ParseFloat(agesSlice[idx], 64)
		if e != nil {
			errs = append(errs, e)
			continue
		}

		job, e := scheduler.NewJob(
			gocron.CronJob(crontab, false),
			gocron.NewTask(func() {
				deleted, errs := fileHandler.DeleteOldFiles(fileSystem, path, age)
				for _, d := range deleted {
					logger.Info("Deleted file", d)
				}

				for _, e := range errs {
					logger.Error("Error deleting file", e.Error())
				}

				if len(deleted) == 0 && len(errs) == 0 {
					logger.Info("No files to delete in path", path)
				}
			}),
			gocron.WithName("PathCleaner-"+path),
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
