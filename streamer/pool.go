package streamer

import "fmt"

// Worker Pool
type VideoDispatcher struct {
	WorkerPool chan chan VideoProcessingJob // channel of channels, enables 2 way communication within a channel
	maxWorkers int
	jobQueue   chan VideoProcessingJob // Send things to our worker pool to process them
	Processor  Processor               // Adapter allows us process the videos
}

// type videoWorker -> this is one of the individual workers in the pool
type videoWorker struct {
	id         int
	jobQueue   chan VideoProcessingJob
	workerPool chan chan VideoProcessingJob // bidirectional channel (https://tleyden.github.io/blog/2013/11/23/understanding-chan-chans-in-go/)
}

// newVideoWorker
func newVideoWorker(id int, workerPool chan chan VideoProcessingJob) videoWorker {
	fmt.Println("newVideoWorker: Creating video worker id", id)
	return videoWorker{
		id:         id,
		jobQueue:   make(chan VideoProcessingJob),
		workerPool: workerPool,
	}
}

// start()
// Anytime start is called it calls an individual worker as a goroutine which executes forever
func (w videoWorker) start() {
	fmt.Println("w.Start(): Starting worker id", w.id)
	go func() {
		for {
			// Add jobQueue to the worker pool
			w.workerPool <- w.jobQueue // whats going on here?

			// Wait for a job to come back (because this go routine will block until something comes in to populate this variable "job")
			job := <-w.jobQueue

			// Process the job
			w.processVideoJob(job.Video)
		}
	}()
}

// run()
func (vd *VideoDispatcher) Run() {
	fmt.Println("vd.Run(): Starting worker pool by running workers")
	for i := 0; i < vd.maxWorkers; i++ {
		fmt.Println("vd.Run(): starting worker id", i+1)
		worker := newVideoWorker(i+1, vd.WorkerPool)
		worker.start()
	}

	go vd.dispatch()
}

// dispatch() (dispatch a worker, assign it a worker)
func (vd *VideoDispatcher) dispatch() {
	for {
		// Wait for a job to come in
		job := <-vd.jobQueue

		fmt.Println("vd.dispatch(): sending", job.Video.ID, "to worker job queue")

		go func() {
			workerJobQueue := <-vd.WorkerPool
			workerJobQueue <- job
		}()
	}
}

// processVideoJob
func (w *videoWorker) processVideoJob(video Video) {
	fmt.Println("w.processVideoJob(): staring encode on video", video.ID)
	video.encode()
}
