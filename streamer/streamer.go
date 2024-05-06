package streamer

type ProcessingMessage struct {
	ID         int
	Successful bool
	Message    string
	OutputFile string
}

// This will hold the unit of work that we want our worker pool to perform
type VideoProcessingJob struct {
	Video Video
}

// This will return the format of the data we need (ex. convert mp4 into a web mp4)
type Processor struct {
}

type Video struct {
	ID         int
	InputFile  string // the video we want to encode
	OutputDir  string // where we want the encoded video to show up
	Type       string
	NotifyChan chan ProcessingMessage // Where are we going to send the processed video to
	// Options *VideoOptions
	Encoder Processor
}

func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	workerPool := make(chan chan VideoProcessingJob, maxWorkers)

	// Todo: implement processor logic
	p := Processor{}

	return &VideoDispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		WorkerPool: workerPool,
		Processor:  p,
	}
}
