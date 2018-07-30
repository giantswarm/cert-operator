package microerror

import (
	"encoding/json"
)

// Error is a predefined error structure whose purpose is to act as container
// for meta information associated to a specific error. The specific error type
// matching can be used as usual. The usual error masking and cause gathering
// can be used as usual. Using Error might look as follows. In the beginning is
// a usual error defined, along with its matcher. This error is the root cause
// once emitted during runtime.
//
//     var notEnoughWorkersError = microerror.Error{
//         Desc: "The amount of requested guest cluster workers exceeds the available number of host cluster nodes.",
//         Docs: "https://github.com/giantswarm/ops-recipes/blob/master/349-not-enough-workers.md",
//         Kind: "notEnoughWorkersError",
//     }
//
//     func IsNotEnoughWorkers(err error) bool {
//         return microerror.Cause(err) == notEnoughWorkersError
//     }
//
type Error struct {
	Desc string `json:"desc"`
	Docs string `json:"docs"`
	Kind string `json:"kind"`
}

func (e *Error) Error() string {
	return e.Desc
}

func (e *Error) GoString() string {
	return e.String()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return nil, Mask(err)
	}

	return b, nil
}

func (e *Error) String() string {
	b, err := e.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (e *Error) UnmarshalJSON(b []byte) error {
	var c Error
	err := json.Unmarshal(b, &c)
	if err != nil {
		return Mask(err)
	}

	*e = c

	return nil
}
