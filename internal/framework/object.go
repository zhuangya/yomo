package framework

import (
	"io"

	txtkv "github.com/10cella/yomo-txtkv-codec"

	"github.com/yomorun/yomo/pkg/plugin"
)

type YomoObjectPluginStream struct {
	Writer YomoObjectPluginStreamWriter
	Reader YomoObjectPluginStreamReader
}

type YomoObjectPluginStreamWriter struct {
	Name   string
	Plugin plugin.YomoObjectPlugin
	io.Writer
}

type YomoObjectPluginStreamReader struct {
	Name string
	io.Reader
}

func (w YomoObjectPluginStreamWriter) Write(b []byte) (int, error) {
	head := b[:1]
	var err error = nil

	var value interface{}
	value, err = txtkv.ObjectCodec{}.Unmarshal(b[1:])
	if err != nil {
		return 0, err
	}

	value, err = w.Plugin.Handle(value) // nolint

	var result []byte
	result, err = txtkv.ObjectCodec{}.Marshal(value.(string)) // nolint

	result = append(head, result...)
	_, err = w.Writer.Write(result) // nolint

	return len(b), err
}

func (r YomoObjectPluginStreamReader) Read(b []byte) (int, error) {
	return r.Reader.Read(b)
}

func NewObjectPlugin(h plugin.YomoObjectPlugin) YomoObjectPluginStream {
	name := "plugin"
	reader, writer := io.Pipe()
	w := YomoObjectPluginStreamWriter{name, h, writer}
	r := YomoObjectPluginStreamReader{name, reader}
	s := YomoObjectPluginStream{Writer: w, Reader: r}
	return s
}
