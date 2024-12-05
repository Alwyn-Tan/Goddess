package parser

import (
	"Goddess/msg"
	"bufio"
	"bytes"
	"io"
)

type Payload struct {
	Data  msg.Msg
	Error error
}

func ParseInputStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse(reader, ch)
	return ch
}

func parse(rawReader io.Reader, ch chan<- *Payload) {
	reader := bufio.NewReader(rawReader)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				close(ch)
				return
			} else {
				ch <- &Payload{
					Error: err,
				}
				close(ch)
				return
			}
		}
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		//Simple Strings
		case '+':
			//Simple Errors
		case '-':
			//Integers
		case ':':
			//Bulk Strings
		case '$':
			//Arrays
		case '*':
		default:
		}
	}
}
