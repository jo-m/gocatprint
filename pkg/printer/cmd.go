package printer

import "strconv"

var (
	cmdGetDevState      = bs([]int8{81, 120, -93, 0, 1, 0, 0, 0, -1})
	cmdSetQuality200DPI = bs([]int8{81, 120, -92, 0, 1, 0, 50, -98, -1})
	cmdLatticeStart     = bs([]int8{81, 120, -90, 0, 11, 0, -86, 85, 23, 56, 68, 95, 95, 95, 68, 56, 44, -95, -1})
	cmdLatticeEnd       = bs([]int8{81, 120, -90, 0, 11, 0, -86, 85, 23, 0, 0, 0, 0, 0, 0, 0, 23, 17, -1})
	cmdSetPaper         = bs([]int8{81, 120, -95, 0, 2, 0, 48, 0, -7, -1})
	cmdPrintText        = bs([]int8{81, 120, -66, 0, 1, 0, 1, 7, -1})
	checksumTable       = bs([]int8{0, 7, 14, 9, 28, 27, 18, 21, 56, 63, 54, 49, 36, 35, 42, 45, 112, 119, 126, 121, 108, 107, 98, 101, 72, 79, 70, 65, 84, 83, 90, 93, -32, -25, -18, -23, -4, -5, -14, -11, -40, -33, -42, -47, -60, -61, -54, -51, -112, -105, -98, -103, -116, -117, -126, -123, -88, -81, -90, -95, -76, -77, -70, -67, -57, -64, -55, -50, -37, -36, -43, -46, -1, -8, -15, -10, -29, -28, -19, -22, -73, -80, -71, -66, -85, -84, -91, -94, -113, -120, -127, -122, -109, -108, -99, -102, 39, 32, 41, 46, 59, 60, 53, 50, 31, 24, 17, 22, 3, 4, 13, 10, 87, 80, 89, 94, 75, 76, 69, 66, 111, 104, 97, 102, 115, 116, 125, 122, -119, -114, -121, -128, -107, -110, -101, -100, -79, -74, -65, -72, -83, -86, -93, -92, -7, -2, -9, -16, -27, -30, -21, -20, -63, -58, -49, -56, -35, -38, -45, -44, 105, 110, 103, 96, 117, 114, 123, 124, 81, 86, 95, 88, 77, 74, 67, 68, 25, 30, 23, 16, 5, 2, 11, 12, 33, 38, 47, 40, 61, 58, 51, 52, 78, 73, 64, 71, 82, 85, 92, 91, 118, 113, 120, 127, 106, 109, 100, 99, 62, 57, 48, 55, 34, 37, 44, 43, 6, 1, 8, 15, 26, 29, 20, 19, -82, -87, -96, -89, -78, -75, -68, -69, -106, -111, -104, -97, -118, -115, -124, -125, -34, -39, -48, -41, -62, -59, -52, -53, -26, -31, -24, -17, -6, -3, -12, -13})
)

// to easily write tests from Python output (bytearray)
func unquote(s string) []byte {
	unq, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		panic(err)
	}

	return []byte(unq)
}

// bs casts signed (Java) bytes ([]int8) to Go signed []byte
func bs(b []int8) (ret []byte) {
	ret = make([]byte, len(b))
	for i, v := range b {
		ret[i] = byte(v)
	}
	return
}

func checkSum(data []byte, startIx, len int) (ret byte) {
	ret = 0
	for i := startIx; i < startIx+len; i++ {
		ret = checksumTable[(ret^data[i])&0xff]
	}

	return
}

func cmdFeedPaper(howMuch int8) (ret []byte) {
	ret = bs([]int8{
		81,
		120,
		-67,
		0,
		1,
		0,
		howMuch,
		0,
		-1,
	})
	ret[7] = checkSum(ret, 6, 1)
	return
}

func cmdSetEnergy(val int16) (ret []byte) {
	ret = bs([]int8{
		81,
		120,
		-81,
		0,
		2,
		0,
		int8((val >> 8) & -1),
		int8(val),
		0,
		-1,
	})
	ret[7] = checkSum(ret, 6, 2)
	return
}
