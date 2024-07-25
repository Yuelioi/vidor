package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"path/filepath"
	"strings"

	"github.com/Yuelioi/vidor/shared"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type LogAdapter struct {
	file *os.File
}

func (l *LogAdapter) Write(p []byte) (n int, err error) {
	return l.file.Write(p)
}

func createLogAdapter(part *shared.Part) (*LogAdapter, error) {
	logFile, err := os.Create(filepath.Join(part.DownloadDir, "ffmpeg.log"))
	if err != nil {
		return nil, err
	}

	return &LogAdapter{file: logFile}, nil
}

// 设置FFmpeg完整路径
func SetFFmpegPath(ffmpegPath string) error {
	if _, err := os.Stat(ffmpegPath); err != nil {
		return err
	}

	err := os.Setenv("FFMPEG_BIN", ffmpegPath)
	if err != nil {
		return err
	}
	if !isFFmpegExecutable(ffmpegPath) {
		return fmt.Errorf("%s is not an ffmpeg executable", ffmpegPath)
	}
	return nil
}

func isFFmpegExecutable(path string) bool {
	cmd := exec.Command(path, "-version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	if !strings.Contains(string(output), "ffmpeg") {
		return false
	}
	return true
}

// func quotePath(path string) string {
// 	return fmt.Sprintf(`"%s"`, path)
// }

// 合并音频与视频
func CombineAV(ffmpegPath string, part *shared.Part, stopChannel chan struct{}) (err error) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-stopChannel:
				fmt.Println("ffmpeg stop")
				cancel()
				return
			case <-ctx.Done():
				fmt.Println("ffmpeg done")
				return
			}
		}
	}()

	input_v := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp4", part.Title))
	input_a := filepath.Join(part.DownloadDir, fmt.Sprintf("%s_temp.mp3", part.Title))
	output_v := filepath.Join(part.DownloadDir, fmt.Sprintf("%s.mp4", part.Title))

	input := []*ffmpeg_go.Stream{ffmpeg_go.Input(input_v), ffmpeg_go.Input(input_a)}
	out := ffmpeg_go.OutputContext(ctx, input, output_v, ffmpeg_go.KwArgs{"c:v": "copy", "c:a": "aac"})

	_, err = os.Stat(ffmpegPath)

	if err == nil {
		out = out.SetFfmpegPath(ffmpegPath)
	}

	// err = out.OverWriteOutput().WithOutput().Run()
	logAdapter, err := createLogAdapter(part)
	if err != nil {
		return err
	}
	defer logAdapter.file.Close()

	cmd := out.OverWriteOutput().WithOutput(logAdapter, logAdapter).Compile()

	// TODO关闭cmd弹窗
	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	return err
}

// 合并音频片段 {"v": 0, "a": 1} {"v": 1, "a": 0}
func CombineSegments(segs []string, out string, args map[string]interface{}) {
	if len(segs) == 0 {
		return
	}

	Inputs := make([]*ffmpeg_go.Stream, len(segs))
	for i, input := range segs {
		Inputs[i] = ffmpeg_go.Input(input)
	}
	err := ffmpeg_go.Concat(Inputs, ffmpeg_go.KwArgs(args)).
		Output(out).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
	if err != nil {
		log.Printf("合并音频文件时出错: %v", err)
	}
}
