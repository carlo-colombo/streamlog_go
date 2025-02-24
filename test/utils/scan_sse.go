package utils

import "bufio"

func ScanEvent(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, bufio.ErrFinalToken
	}
	advance, token, err = bufio.ScanLines(data, false)
	if data[advance] == '\n' {
		return advance + 1, token, err
	} else {
		nAdvance, nToken, _ := ScanEvent(data[advance:], false)

		merged := append(token, '\n')
		merged = append(merged, nToken...)
		return advance + nAdvance, merged, nil
	}
}
