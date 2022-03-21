package ffmpeg

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"
)

func LogOutput(r io.Reader, done <-chan bool) {
	log.Println("output logger started")
	s := bufio.NewScanner(r)
loop:
	for {
		select {
		case <-done:
			break loop
		default:
			s.Scan()
			fmt.Println(s.Text())
			if s.Err() == io.EOF {
				fmt.Println("eof")
				break loop
			}
			fmt.Println("wait for output")
			time.Sleep(time.Millisecond * 500)
		}
	}
	log.Println("output logger done")
}
