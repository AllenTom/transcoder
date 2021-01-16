package ffmpeg

import (
	"os/exec"
	"strings"
)

type EncoderFlag struct {
	Video                    bool
	Audio                    bool
	Subtitle                 bool
	FrameLevelMultithreading bool
	SliceLevelMultithreading bool
	Experimental             bool
	DrawHorizBand            bool
	DirectRenderingMethod1   bool
}
type Encoder struct {
	Flags EncoderFlag
	Name  string
	Desc  string
}

func ReadEncoderList(config *Config) ([]Encoder, error) {
	cmd := exec.Command(config.FfmpegBinPath, "-encoders")
	outputByte, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	startFlag := false
	encoderList := make([]Encoder, 0)
	for _, line := range strings.Split(string(outputByte), "\n") {
		line = strings.TrimSpace(line)
		if startFlag {
			encoder := Encoder{
				Flags: EncoderFlag{},
			}
			for idx, part := range strings.SplitN(line, " ", 3) {
				part = strings.TrimSpace(part)
				if idx == 0 {
					for fIdx, flagS := range strings.Split(part, "") {
						if fIdx == 0 {
							switch flagS {
							case "V":
								encoder.Flags.Video = true
								break
							case "A":
								encoder.Flags.Audio = true
								break
							case "S":
								encoder.Flags.Subtitle = true
								break
							}
						}
						if fIdx == 1 && flagS == "F" {
							encoder.Flags.FrameLevelMultithreading = true
						}
						if fIdx == 2 && flagS == "S" {
							encoder.Flags.SliceLevelMultithreading = true
						}
						if fIdx == 3 && flagS == "X" {
							encoder.Flags.Experimental = true
						}
						if fIdx == 4 && flagS == "B" {
							encoder.Flags.DrawHorizBand = true
						}
						if fIdx == 5 && flagS == "D" {
							encoder.Flags.DirectRenderingMethod1 = true
						}
					}
				}
				if idx == 1 {
					encoder.Name = part
				}
				if idx == 2 {
					encoder.Desc = part
				}

			}
			encoderList = append(encoderList, encoder)
		}
		if line == "------" {
			startFlag = true
		}
	}
	return encoderList, err
}
