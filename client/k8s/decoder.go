package k8s

import (
	"encoding/json"
	"io"

	"github.com/giantswarm/certificatetpr"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
)

type CertificateDecoder struct {
	Stream io.ReadCloser
}

// Decode deserializes a runtime object into a certificatetpr CustomObject.
func (d *CertificateDecoder) Decode() (action watch.EventType, object runtime.Object, err error) {
	decoder := json.NewDecoder(d.Stream)

	var e struct {
		Type   watch.EventType
		Object certificatetpr.CustomObject
	}
	if err := decoder.Decode(&e); err != nil {
		return watch.Error, nil, err
	}
	return e.Type, &e.Object, nil
}

func (d *CertificateDecoder) Close() {
	d.Stream.Close()
}
