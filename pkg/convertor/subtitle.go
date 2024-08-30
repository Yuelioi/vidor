package convertor

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func TranscriptToSrt(transcript youtube.VideoTranscript) string {
	srt := ""

	for index, seg := range transcript {
		srt += generateSRTBlock(seg, index+1)
	}
	return srt
}

func WriteSrt(filepath string, srt string) {

	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := strings.NewReader(srt)

	_, err = io.Copy(file, reader)
	if err != nil {
		panic(err)
	}
}

func generateSRTBlock(seg youtube.TranscriptSegment, index int) string {
	return fmt.Sprintf("%d\n%s\n%s\n\n", index, msToSRTTimeRange(seg), seg.Text)
}

func msToTime(ms int) string {
	hours := ms / (1000 * 60 * 60)
	minutes := (ms % (1000 * 60 * 60)) / (1000 * 60)
	seconds := (ms % (1000 * 60)) / 1000
	milliseconds := ms % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}

func msToSRTTimeRange(seg youtube.TranscriptSegment) string {
	return msToTime(seg.StartMs) + " --> " + msToTime(seg.StartMs+seg.Duration)
}
