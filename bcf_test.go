package bcf

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestBcf(t *testing.T) {
	rdr, err := os.Open("tests/dbsnp.sub.bcf")
	//rdr, err := os.Open("u.bcf")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewReader(rdr, 1)
	if err != nil {
		t.Fatal(err)
	}
	for {
		variant, err := b.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(variant.Id), variant.Chrom, variant.Pos)
		//fmt.Println(variant.Info().Get("DP"))
	}
}
