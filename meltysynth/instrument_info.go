package meltysynth

import (
	"encoding/binary"
	"errors"
	"io"
)

type instrumentInfo struct {
	name           string
	zoneStartIndex int32
	zoneEndIndex   int32
}

func readInstrumentsFromChunk(reader io.Reader, size int32) ([]*instrumentInfo, error) {

	var err error

	if size%22 != 0 {
		return nil, errors.New("The instrument list is invalid.")
	}

	count := size / 22

	instruments := make([]*instrumentInfo, count, count)

	for i := int32(0); i < count; i++ {

		instrument := new(instrumentInfo)

		instrument.name, err = readFixedLengthString(reader, 20)
		if err != nil {
			return nil, err
		}

		var zoneStartIndex uint16
		err = binary.Read(reader, binary.LittleEndian, &zoneStartIndex)
		if err != nil {
			return nil, err
		}
		instrument.zoneStartIndex = int32(zoneStartIndex)

		instruments[i] = instrument
	}

	for i := int32(0); i < count-1; i++ {
		instruments[i].zoneEndIndex = instruments[i+1].zoneStartIndex - 1
	}

	return instruments, nil
}