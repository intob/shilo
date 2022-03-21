package ffmpeg

import (
	"fmt"
	"os/exec"
)

func Scale(in, out, outRes string) *exec.Cmd {
	scale := fmt.Sprintf("scale=%s", outRes)
	cmd := exec.Command("ffmpeg",
		"-i", in,
		"-vf", scale,
		out)
	return cmd
}
