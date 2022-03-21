package ffmpeg

import (
	"context"
	"os/exec"
)

func Scale(ctx context.Context, in, out, outRes string) *exec.Cmd {
	scale := "scale=" + outRes
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", in,
		"-vf", scale,
		out)
	return cmd
}
