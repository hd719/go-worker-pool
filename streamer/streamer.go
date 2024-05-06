package streamer

type ProcessingMessage struct {
	ID         int
	Successful bool
	Message    string
	OutputFile string
}

// This will hold the unit of work that we want our worker pool to perform
// We wrap this type around a Video, which has all the information we need about the input source and what we want the output to look like()
type VideoProcessingJob struct {
	Video Video
}

// This will return the format of the data we need (ex. convert mp4 into a web mp4)
// Does the actually processing
type Processor struct {
	Engine Encoder
}

type Video struct {
	ID           int
	InputFile    string // the video we want to encode
	OutputDir    string // where we want the encoded video to show up
	Type         string
	NotifyChan   chan ProcessingMessage // Where are we going to send the processed video to
	Options      *VideoOptions
	Encoder      Processor
	EncodingType string
}

type VideoOptions struct {
	RenameOutput    bool
	SegmentDuration int
	MaxRate1080p    string
	MaxRate720p     string
	MaxRate480p     string
}

func (vd *VideoDispatcher) NewVideo(id int, input string, output string, encType string, notifyChan chan ProcessingMessage, options *VideoOptions) Video {
	if options == nil {
		options = &VideoOptions{}
	}

	return Video{
		ID:           id,
		InputFile:    input,
		OutputDir:    output,
		EncodingType: encType,
		NotifyChan:   notifyChan,
		Encoder:      vd.Processor,
		Options:      options,
	}
}

func (v *Video) encode() {
}

// New creates and returns a new worker pool
func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	workerPool := make(chan chan VideoProcessingJob, maxWorkers)

	// Todo: implement processor logic
	var e VideoEncoder
	p := Processor{
		Engine: &e,
	}

	return &VideoDispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		WorkerPool: workerPool,
		Processor:  p,
	}
}
