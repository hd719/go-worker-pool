package streamer

// Worker Pool
type VideoDispatcher struct {
	WorkerPool chan chan VideoProcessingJob // channel of channels, enables 2 way communication within a channel
	maxWorkers int
	jobQueue   chan VideoProcessingJob // Send things to our worker pool to process them
	Processor  Processor // Adapter allows us process the videos
}
