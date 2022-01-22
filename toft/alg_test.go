package toft

import (
	"testing"
)

func TestRandomCreateZHCNUnicode(t *testing.T) {
	as := make(map[string]interface{})
	count := 0
	degrees := 1000000
	for i := 0; i < degrees; i++ {
		uStr, str := RandomCreateZHCNUnicode()
		if _, ok := as[uStr]; !ok {
			as[uStr] = str
		} else {
			count++
			//t.Logf("重复,[%s][%s]\n", uStr, str)
		}
	}

	t.Logf("合计重复,[%f%%]\n", float64(count)/float64(degrees))

}

func BenchmarkRandomCreateZHCNUnicode(b *testing.B) {
	as := make(map[string]interface{})
	count := 0
	degrees := 1000000
	for i := 0; i < degrees; i++ {
		uStr, str := RandomCreateZHCNUnicode()
		if _, ok := as[uStr]; !ok {
			as[uStr] = str
		} else {
			count++
			//t.Logf("重复,[%s][%s]\n", uStr, str)
		}
	}

	b.Logf("合计重复,[%f%%]\n", float64(count)/float64(degrees))
}
