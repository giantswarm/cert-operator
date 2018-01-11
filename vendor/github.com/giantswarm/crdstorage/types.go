package crdstorage

type storageConfigJSONPatch struct {
	Spec storageConfigJSONPatchSpec `json:"spec"`
}

type storageConfigJSONPatchSpec struct {
	Storage storageConfigJSONPatchSpecStorage `json:"storage" yaml:"storage"`
}

type storageConfigJSONPatchSpecStorage struct {
	Data map[string]*string `json:"data" yaml:"data"`
}
