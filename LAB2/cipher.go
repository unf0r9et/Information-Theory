package main

const Degree = 25

var Taps = []int{25, 3}

type LFSR struct {
	state []byte
	taps  []int
}

func NewLFSR(seed string, taps []int) *LFSR {
	state := make([]byte, len(seed))
	for i, char := range seed {
		if char == '1' {
			state[i] = 1
		} else {
			state[i] = 0
		}
	}
	return &LFSR{state: state, taps: taps}
}

func (l *LFSR) NextBit() byte {
	outBit := l.state[0]

	feedback := l.state[0] ^ l.state[22]

	for i := 0; i < 24; i++ {
		l.state[i] = l.state[i+1]
	}

	l.state[24] = feedback

	return outBit
}

func (l *LFSR) ProcessData(inputData []byte, limitBytes int) (outputData []byte, firstKey, lastKey []byte) {
	outputData = make([]byte, len(inputData))
	firstKey = make([]byte, 0, limitBytes)
	lastKeyRing := make([]byte, limitBytes)

	for i, b := range inputData {
		var cryptByte byte = 0
		var keyByte byte = 0

		for bitIdx := 7; bitIdx >= 0; bitIdx-- {
			keyBit := l.NextBit()
			fileBit := (b >> bitIdx) & 1
			cryptBit := fileBit ^ keyBit

			keyByte |= (keyBit << bitIdx)
			cryptByte |= (cryptBit << bitIdx)
		}

		outputData[i] = cryptByte

		if i < limitBytes {
			firstKey = append(firstKey, keyByte)
		}
		if len(inputData) > 0 {
			lastKeyRing[i%limitBytes] = keyByte
		}
	}

	if len(inputData) <= limitBytes {
		lastKey = firstKey
	} else {
		lastKey = make([]byte, limitBytes)
		for i := 0; i < limitBytes; i++ {
			lastKey[i] = lastKeyRing[(len(inputData)+i)%limitBytes]
		}
	}

	return outputData, firstKey, lastKey
}
