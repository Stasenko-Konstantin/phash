based on [phash algorithm](https://www.phash.org/docs/pubs/thesis_zauner.pdf)

```go
import "github.com/Stasenko-Konstantin/phash"

hash1, err := phash.Hash("path/to/img1")
if err != nil { ... }
hash2, err := phash.Hash("path/to/img2")
if err != nil { ... }
if phash.Distance(hash1, hash2) > phash.MinDistance {
	... // rm imgs or smth else
}
```