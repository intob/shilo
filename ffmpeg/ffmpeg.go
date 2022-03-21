package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"
)

func Sanitise(input string) string {
	c := strings.ReplaceAll(input, ";", "")
	return strings.ReplaceAll(c, "&&", "")
}

func Scale(in, out string, width, height int) *exec.Cmd {
	scale := fmt.Sprintf("scale=%v:%v", width, height)
	cmd := exec.Command("ffmpeg",
		"-i", in,
		"-vf", scale,
		out)
	fmt.Println(cmd.String())
	return cmd
}
