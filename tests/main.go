package main

import (
	"fmt"
	"io"
	"os"

	"github.com/brentp/bcf"
)

func main() {
	rdr, _ := os.Open("tests/dbsnp.sub.bcf")
	brdr, _ := bcf.NewReader(rdr, 2)
	fmt.Println(string(brdr.Header.Text))
	for {
		variant, err := brdr.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(variant.Chrom, variant.Id, variant.Pos, variant.Info())
	}
}
