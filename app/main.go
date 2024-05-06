package main

import (
	"fmt"
	"streamer"
)

func main() {
	// Define number of workers and jobs
	const numJobs = 1
	const numWorkers = 2

	// Create 2 channels for work and results (1. notifications and 2. We send work to)
	notifyChan := make(chan streamer.ProcessingMessage, numJobs)
	defer close(notifyChan)

	videoQueue := make(chan streamer.VideoProcessingJob, numJobs)
	defer close(videoQueue)

	// Get a worker pool
	wp := streamer.New(videoQueue, numWorkers)

	// Start the worker pool
	wp.Run()
	fmt.Println("Worker Pool Started. Press enter to continue")
	_, _ = fmt.Scanln()

	// Create 1 video to send to the worker pool
	ops := &streamer.VideoOptions{
		RenameOutput:    true,
		SegmentDuration: 10,
		MaxRate1080p:    "1200k",
		MaxRate720p:     "600k",
		MaxRate480p:     "400k",
	}
	video := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "hls", notifyChan, ops)
	// video := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "mp4", notifyChan, ops)

	// Send the videos to the worker pool
	videoQueue <- streamer.VideoProcessingJob{
		Video: video,
	}

	// Print out the results
	for i := 1; i <= numJobs; i++ {
		msg := <-notifyChan
		fmt.Println("i:", i, "/", "message:", msg)
	}

	fmt.Println("Done")
}
