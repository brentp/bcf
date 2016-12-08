bcf
===

[![Build Status](https://travis-ci.org/brentp/bcf.svg?branch=master)](https://travis-ci.org/brentp/bcf)

`bcf` is a [bcf](https://samtools.github.io/hts-specs/BCFv2_qref.pdf) parser for the go programming language.

```Go
import "github.com/brentp/bcf"

func main() {
    rdr, _ := os.Open("some.bcf")
    brdr, _ := bcf.NewReader(rdr, 2)
    for {
        variant, err := brdr.Read()
        if err == io.EOF {
            break
        }
        fmt.Println(variant.Chrom, variant.Id, variant.Pos)
    }
}
```

TODO
====
The library is currently working, but some things remain to be done:

+ fix parsing of FILTER
+ finalize parsing of FORMAT fields. This is done, but we don't currently replace the missing tokens with NaN or whatever is appropriate.
+ benchmark the current INFO and FORMAT parsing (which is lazy) that saves a slice (3 uint32's) to the underlying data to compare to just saving the offsets and having each entry in the INFO pull from
 the same underlying slice.
+ docs
