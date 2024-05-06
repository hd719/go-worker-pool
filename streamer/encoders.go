package streamer

import (
	"fmt"

	"github.com/xfrr/goffmpeg/transcoder"
)

// Encoder is an interface for encoding video, any type that wants to satisfy this interface must implement all its methods
type Encoder interface {
	EncodeToMP4(v *Video, baseFileName string) error
}

// VideoEncoder is a type which satisfies the Encoder interface because it implements all the methods specified in Encoder
type VideoEncoder struct{}

// Takes a video object and a base file name and encodes to mp4
func (ve *VideoEncoder) EncodeToMP4(v *Video, baseFileName string) error {
	// Create transcoder
	trans := new(transcoder.Transcoder)

	// Build the output path
	outputPath := fmt.Sprintf("%s/%s", v.OutputDir, baseFileName)

	// Initialize the transcoder
	err := trans.Initialize(v.InputFile, outputPath)
	if err != nil {
		fmt.Println("Error encoding....")
		return err
	}

	// Set Codec.
	trans.MediaFile().SetVideoCodec("libx264")

	// Start the transcoding process
	done := trans.Run(false)

	err = <-done
	if err != nil {
		return err
	}

	return nil
}
