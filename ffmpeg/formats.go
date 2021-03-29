package ffmpeg

import (
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

func GetFormats(config *Config) ([]SupportFormat, error) {
	cmd := exec.Command(config.FfmpegBinPath, "-formats")
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
			if strings.HasPrefix(flagText, "D") {
				format.Flags.Demuxing = true
			} else if strings.HasPrefix(flagText, "E") {
				format.Flags.Muxing = true
			} else if strings.HasPrefix(flagText, "DE") {
				format.Flags.Muxing = true
				format.Flags.Demuxing = true
			}
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "D ") {
				line = strings.Replace(line, "D ", "D", 1)
			}
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
			formats = append(formats, format)
		}
		if strings.TrimSpace(line) == "--" {
			startFlag = true
		}
	}
	return formats, nil
}
