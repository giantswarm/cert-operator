package crt

import (
	"encoding/json"
	microerror "github.com/giantswarm/microkit/error"
	"io"

	"github.com/giantswarm/certificatetpr"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
)

// certificateDecoder is used by the watcher to decode the certificatetpr objects.
type certificateDecoder struct {
	decoder *json.Decoder
	close   func() error
}

// newCertificateDecoder returns the decoder.
func newCertificateDecoder(stream io.ReadCloser) *certificateDecoder {
	return &certificateDecoder{
		decoder: json.NewDecoder(stream),
		close:   stream.Close,
	}
}

// Decode deserializes a runtime object into a certificatetpr CustomObject.
func (cd *certificateDecoder) Decode() (watch.EventType, runtime.Object, error) {
	var e struct {
		Type   watch.EventType
		Object certificatetpr.CustomObject
	}

	if err := cd.decoder.Decode(&e); err != nil {
		return watch.Error, nil, microerror.MaskAny(err)
	}

	return e.Type, &e.Object, nil
}

// Close closes the stream.
func (cd *certificateDecoder) Close() {
	err := cd.close()
	if err != nil {
		panic(err)
	}
}
