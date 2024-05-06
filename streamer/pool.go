package streamer

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
	return videoWorker{
		id:         id,
		jobQueue:   make(chan VideoProcessingJob),
		workerPool: workerPool,
	}
}

// start()
// Anytime start is called it calls an individual worker as a goroutine which executes forever
func (w videoWorker) start() {
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
	for i := 0; i < vd.maxWorkers; i++ {
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

		go func() {
			workerJobQueue := <-vd.WorkerPool
			workerJobQueue <- job
		}()
	}
}

// processVideoJob
func (w *videoWorker) processVideoJob(video Video) {
	video.encode()
}
