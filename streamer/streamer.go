package streamer

type ProcessingMessage struct {
	ID         int
	Successful bool
	Message    string
	OutputFile string
}

type VideoProcessingJob struct {
}
