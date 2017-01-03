package obj

import (
	"strconv"
	_ "encoding/hex"
	_ "fmt"
	"encoding/hex"
)

const (
	errInvalidFormat     = "Invalid format"
	errHeadLoaded        = "Head already loaded"
	errEndLoaded         = "End already loaded"
	errInvalidHeadFormat = "Invalid head format"
	errInvalidBodyFormat = "Invalid body format"
	errInvalidEndFormat  = "Invalid end format"
)

type BodyObjectCode struct {
	StartAddr int32
	Length    int32
	Code      []byte
}

// ObjectCode ...
type ObjectCode struct {
	Name      string
	Length    int32
	LoadAddr  int32
	StartAddr int32
	Code      []*BodyObjectCode
	// Hidden
	headLoaded bool
	endLoaded  bool
}

// Load ...
func (obj *ObjectCode) Load(bytes []byte) {
	if len(bytes) > 0 {
		switch bytes[0] {
		case 'H':
			obj.loadHead(bytes)
		case 'T':
			obj.loadBody(bytes)
		case 'E':
			obj.loadEnd(bytes)
		default:
			panic(errInvalidFormat)
		}
	}
}

func (obj *ObjectCode) loadHead(str []byte) {
	if obj.headLoaded {
		panic(errHeadLoaded)
	}

	if len(str) < 19 || str[0] != 'H' {
		panic(errInvalidHeadFormat)
	}

	obj.Name = string(str[1:7])

	if loadAddr64, err := strconv.ParseUint(string(str[7:13]), 16, 64); err == nil {
		obj.LoadAddr = int32(loadAddr64)
	} else {
		panic(err)
	}

	if progLen64, err := strconv.ParseUint(string(str[13:19]), 16, 64); err == nil {
		obj.Length = int32(progLen64)
	} else {
		panic(err)
	}

	obj.headLoaded = true
}
func (obj *ObjectCode) loadBody(str []byte) {
	if len(str) < 10 || str[0] != 'T' {
		panic(errInvalidBodyFormat)
	}

	body := &BodyObjectCode{}

	// Save the start address of this block
	if startAddr64, err := strconv.ParseUint(string(str[1:7]), 16, 64); err != nil {
		panic(err)
	} else {
		body.StartAddr = int32(startAddr64)
	}

	if codeLen, err := strconv.ParseUint(string(str[7:9]), 16, 64); err != nil {
		panic(err)
	} else {
		body.Length = int32(codeLen)
	}

	for i := int32(0); i < 2*body.Length; i += 2 {
		b, _ := hex.DecodeString(string(str[9 + i:9 + i + 2]))
		// Lets just ignore the error
		body.Code = append(body.Code, b[0])
	}

	obj.Code = append(obj.Code, body)
}
func (obj *ObjectCode) loadEnd(str []byte) {
	if obj.endLoaded {
		panic(errEndLoaded)
	}

	if len(str) < 7 || str[0] != 'E' {
		panic(errInvalidEndFormat)
	}

	if progStart64, err := strconv.ParseUint(string(str[1:]), 16, 64); err != nil {
		panic(err)
	} else {
		obj.StartAddr = int32(progStart64)
	}

	obj.endLoaded = true
}
