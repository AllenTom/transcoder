package ffmpeg

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
