# splitzip

Allows you to define a maximum file size, and will split zip file output based
on that size. Each zip file produced is a completely functional archive, i.e. it
doesn't generate a single huge zip file then split that into pieces.

# Example

Given a folder with a set of files of various sizes

```
sources/10mb-32.txt	10000000
sources/11mb-32.txt	11000000
sources/1mb-32.txt	1000000
sources/2mb-32.txt	2000000
sources/3mb-32.txt	3000000
sources/4mb-32.txt	4000000
sources/5mb-32.txt	5000000
sources/6mb-32.txt	6000000
sources/7mb-32.txt	7000000
sources/8mb-32.txt	8000000
sources/9mb-32.txt	9000000
```

the following program will archive all the text files above

```
package main

import (
	"fmt"
	"github.com/splitzip"
)

func main() {
	fmt.Printf("Zip(), any error? %+v\n", splitzip.Zip("sources", "sources/zip", "test", 10000000, .4))
}
```

and will produce a numbered sequence of zip files smaller than the defined
"max size" of `10MB`:

```
sources/zip/test0.zip	9444579
sources/zip/test1.zip	9444588
sources/zip/test2.zip	9444622
sources/zip/test3.zip	9444513
sources/zip/test4.zip	3778674
```

The final parameter `.4` in the call to `splitzip.Zip()` above is the expected
space savings for this set of files. Because the density of files will vary, and
hence their compressibility, you can use this number to optimize for your set of
file types.
