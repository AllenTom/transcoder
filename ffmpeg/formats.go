package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"
)

type FormatFlag struct {
	Demuxing bool
	Muxing   bool
}
type SupportFormat struct {
	Flags FormatFlag
	Name  string
	Desc  string
}

func GetFormats() ([]SupportFormat, error) {
	cmd := exec.Command("ffmpeg", "-formats")
	outputByte, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	formats := make([]SupportFormat, 0)
	startFlag := false
	for _, line := range strings.Split(string(outputByte), "\n") {
		format := SupportFormat{
			Flags: FormatFlag{},
		}
		if startFlag && len(line) > 0 {
			flagText := line[:3]
			flagText = strings.TrimSpace(flagText)
			switch flagText {
			case "D":
				format.Flags.Demuxing = true
				break
			case "E":
				format.Flags.Muxing = true
				break
			case "DE":
				format.Flags.Muxing = true
				format.Flags.Demuxing = true
				break
			}
			line = strings.TrimSpace(line)
			for idx, part := range strings.SplitN(line, " ", 3) {
				part = strings.TrimSpace(part)
				if idx == 0 {
					switch part {
					case "D":
						format.Flags.Demuxing = true
						break
					case "E":
						format.Flags.Muxing = true
						break
					case "DE":
						format.Flags.Muxing = true
						format.Flags.Demuxing = true
						break
					}
				}
				if idx == 1 {
					format.Name = part
				}
				if idx == 2 {
					format.Desc = part
				}
			}
			fmt.Println(format)
			formats = append(formats, format)
		}
		if line == " --" {
			startFlag = true
		}
	}
	return formats, nil
}
