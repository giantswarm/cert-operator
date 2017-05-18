package certificatetpr

import (
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/client-go/pkg/api/unversioned"
)

// List represents a list of CustomObject resources.
type List struct {
	unversioned.TypeMeta `json:",inline"`
	Metadata             unversioned.ListMeta `json:"metadata"`

	Items []CustomObject `json:"items"`
}

// GetObjectKind is required to satisfy the Object interface.
func (l *List) GetObjectKind() unversioned.ObjectKind {
	return &l.TypeMeta
}

// GetListMeta is required to satisfy the ListMetaAccessor interface.
func (l *List) GetListMeta() unversioned.List {
	return &l.Metadata
}

// The code below is used only to work around a known problem with third-party
// resources and ugorji. If/when these issues are resolved, the code below
// should no longer be required.

type ListCopy List

func (l *List) UnmarshalJSON(data []byte) error {
	fmt.Fprintf(os.Stdout, "IN List UnmarshalJSON data: %s \n", data)

	tmp := ListCopy{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "TMP List: #%v \n", tmp)

	tmp2 := List(tmp)
	*l = tmp2
	fmt.Fprintf(os.Stdout, "RETURNED List: #%v \n", tmp2)
	return nil
}
