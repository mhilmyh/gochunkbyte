package gochunkbyte

type chunk struct {
	Bytes []byte `json:"bytes"`
	Index int    `json:"index"`
	Start int    `json:"start"`
}

func Split(content []byte, size int) ([]*chunk, int) {
	l := len(content)
	chunks := make([]*chunk, size)
	n := l / size
	r := l % size
	ch := make(chan *chunk)
	var start, end, counter int
	for i := 0; i < size; i++ {
		start = end
		end = start + n
		if r != 0 {
			end++
			r--
		}
		go makeChunk(ch, content, i, start, end)
		counter++
	}
	for counter != 0 {
		c := <-ch
		if c != nil {
			chunks[c.Index] = c
		}
		counter--
	}
	close(ch)
	return chunks, l
}

func makeChunk(ch chan *chunk, content []byte, index, start, end int) {
	c := chunk{
		Index: index,
		Start: start,
	}
	for i := start; i < end; i++ {
		c.Bytes = append(c.Bytes, content[i])
	}
	ch <- &c
}

func Merge(chunks []*chunk, length int) []byte {
	bytes := make([]byte, length)
	ch := make(chan bool)
	counter := 0
	for _, c := range chunks {
		go combineChunk(ch, bytes, c)
		counter++
	}
	for counter != 0 {
		done := <-ch
		if done {
			counter--
		}
	}
	return bytes
}

func combineChunk(ch chan bool, bytes []byte, c *chunk) {
	end := c.Start + len(c.Bytes)
	for i := c.Start; i < end; i++ {
		bytes[i] = c.Bytes[i-c.Start]
	}
	ch <- true
}
