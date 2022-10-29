package gochunkbyte

import "log"

type chunk struct {
	bytes []byte
	index int
}

func Split(content []byte, size int) []*chunk {
	l := len(content)
	n := l / size
	chunks := make([]*chunk, n)
	ch := make(chan chunk)
	for i := 0; i < l; i += n {
		go makeChunk(ch, content, i, i+n)
	}
	select {
	case c := <-ch:
		log.Println(c)
	}
	close(ch)
	return chunks
}

func makeChunk(ch chan chunk, content []byte, start, end int) {
	c := chunk{}
	for i := start; i < end; i++ {
		c.bytes = append(c.bytes, content[i])
	}
	ch <- c
}
