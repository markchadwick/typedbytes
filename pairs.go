package typedbytes

// ----------------------------------------------------------------------------
// Pair Reader
// ----------------------------------------------------------------------------

type PairReader struct {
	reader *Reader
}

func NewPairReader(reader *Reader) *PairReader {
	return &PairReader{reader}
}

func (pr *PairReader) Next() (k, v interface{}, err error) {
	if k, err = pr.reader.Next(); err != nil {
		return
	}
	if v, err = pr.reader.Next(); err != nil {
		return
	}
	return
}

// ----------------------------------------------------------------------------
// Pair Writer
// ----------------------------------------------------------------------------

type PairWriter struct {
	writer *Writer
}

func NewPairWriter(writer *Writer) *PairWriter {
	return &PairWriter{writer}
}

func (pw *PairWriter) Write(k, v interface{}) (err error) {
	if err = pw.writer.Write(k); err != nil {
		return
	}
	if err = pw.writer.Write(v); err != nil {
		return
	}
	return
}

func (pw *PairWriter) Close() error {
	return pw.writer.Close()
}
