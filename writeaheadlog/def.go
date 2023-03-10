package writeaheadlog

type Record struct {
	Offset uint32 `json:"value"`
	Value  []byte `json:"offset"`
}

var NoRecord = Record{}
