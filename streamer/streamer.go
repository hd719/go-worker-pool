package streamer

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/tsawler/toolbox"
)

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

	fmt.Println("NewVideo() -> New Video created:", id, input)

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

// All pushes to the notify chan will be in this func
func (v *Video) encode() {
	var fileName string

	switch v.EncodingType {
	case "mp4":
		fmt.Println("v.encode(): About to encode to mp4", v.ID)
		// encode the video
		name, err := v.encodeToMp4()
		if err != nil {
			// send info to notify chan
			v.sendToNotifyChan(false, "", fmt.Sprintf("encode failed for %d %s", v.ID, err.Error()))
		}
		fileName = fmt.Sprintf("%s.mp4", name)
	case "hls":
		fmt.Println("v.encode(): About to encode to mp4", v.ID)
		// encode the video
		name, err := v.encodeToHLS()
		if err != nil {
			// send info to notify chan
			v.sendToNotifyChan(false, "", fmt.Sprintf("encode failed for %d %s", v.ID, err.Error()))
		}
		fileName = fmt.Sprintf("%s.m3u8", name)

	default:
		fmt.Println("v.encode(): error trying to encode video", v.ID)
		v.sendToNotifyChan(false, "", fmt.Sprintf("error processing for %d: invalid encoding type", v.ID))
		return
	}

	fmt.Println("v.encode(): sending success message for video id", v.ID, "to notify chan")
	v.sendToNotifyChan(true, fileName, fmt.Sprintf("video id %d processed and saved as %s", v.ID, fmt.Sprintf("%s/%s", v.OutputDir, fileName)))
}

func (v *Video) encodeToMp4() (string, error) {
	baseFileName := ""
	fmt.Println("v.encodeToMP4: about to try to encode video id", v.ID)

	if !v.Options.RenameOutput {
		// Get the base filename
		b := path.Base(v.InputFile)
		baseFileName = strings.TrimSuffix(b, filepath.Ext(b)) // ex. cat.mp4 becomes cat
	} else {
		var t toolbox.Tools
		baseFileName = t.RandomString(10)
	}

	// Encode
	err := v.Encoder.Engine.EncodeToMP4(v, baseFileName)
	if err != nil {
		return "", err
	}
	fmt.Println("v.encodeToMp4: successfully encoded video id", v.ID)

	return baseFileName, nil
}

func (v *Video) encodeToHLS() (string, error) {
	baseFileName := ""

	if !v.Options.RenameOutput {
		// Get the base filename
		b := path.Base(v.InputFile)
		baseFileName = strings.TrimSuffix(b, filepath.Ext(b)) // ex. cat.mp4 becomes cat
	} else {
		var t toolbox.Tools
		baseFileName = t.RandomString(10)
	}

	// Encode
	err := v.Encoder.Engine.EncodeToHLS(v, baseFileName)
	if err != nil {
		return "", err
	}
	fmt.Println("v.encodeToHLS: successfully encoded video id", v.ID)

	return baseFileName, nil
}

func (v *Video) sendToNotifyChan(successful bool, fileName, message string) {
	fmt.Println("v.sendToNotifyChan(): sending message to notifyChan for video id", v.ID)
	v.NotifyChan <- ProcessingMessage{
		ID:         v.ID,
		Successful: successful,
		Message:    message,
		OutputFile: fileName,
	}
}

// New creates and returns a new worker pool
func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	fmt.Println("New: Creating worker pool")
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
