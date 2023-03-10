package writeaheadlog

type Record struct {
	Offset uint32 `json:"offset"`
	Value  []byte `json:"value"`
}

var NoRecord = Record{}
