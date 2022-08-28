package randomstring

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math"
)

/*
func main() {
	for i := 0; i < 50; i++ {
		s := RandomBase64String(i)
		fmt.Printf("%d: %s %t\n", i, s, len(s) == i)
	}

	for i := 0; i < 50; i++ {
		s := RandomBase16String(i)
		fmt.Printf("%d: %s %t\n", i, s, len(s) == i)
	}

}
*/

func RandomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	// fmt.Printf("%v", buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // strip the one extra byte we get from half the results.
}

func RandomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}
