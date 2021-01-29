package ffmpeg

import (
	"os/exec"
	"strings"
)

type CodecFlag struct {
	Decoding            bool
	Encoding            bool
	VideoCodec          bool
	AudioCodec          bool
	SubtitleCodec       bool
	IntraFrameOnly      bool
	LossyCompression    bool
	LosslessCompression bool
}

type Codec struct {
	Flags CodecFlag
	Name  string
	Desc  string
}

func ReadCodecList(config *Config) ([]Codec, error) {
	cmd := exec.Command(config.FfmpegBinPath, "-codecs")
	outputByte, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	startFlag := false
	codecList := make([]Codec, 0)
	for _, line := range strings.Split(string(outputByte), "\n") {
		if startFlag && len(line) > 0 {
			codec := Codec{
				Flags: CodecFlag{},
			}
			line = strings.TrimSpace(line)
			for idx, part := range strings.SplitN(line, " ", 3) {
				part = strings.TrimSpace(part)
				if idx == 0 {
					for fIdx, flagS := range strings.Split(part, "") {
						if fIdx == 0 && flagS == "D" {
							codec.Flags.Decoding = true
						}
						if fIdx == 1 && flagS == "E" {
							codec.Flags.Encoding = true
						}
						if fIdx == 2 {
							switch flagS {
							case "V":
								codec.Flags.VideoCodec = true
								break
							case "A":
								codec.Flags.AudioCodec = true
								break
							case "S":
								codec.Flags.SubtitleCodec = true
								break
							}
						}
						if fIdx == 3 && flagS == "I" {
							codec.Flags.IntraFrameOnly = true
						}
						if fIdx == 4 && flagS == "L" {
							codec.Flags.LossyCompression = true
						}
						if fIdx == 5 && flagS == "S" {
							codec.Flags.LosslessCompression = true
						}
					}
				}
				if idx == 1 {
					codec.Name = part
				}
				if idx == 2 {
					codec.Desc = part
				}

			}
			codecList = append(codecList, codec)
		}
		if strings.TrimSpace(line) == "-------" {
			startFlag = true
		}
	}
	return codecList, nil
}
