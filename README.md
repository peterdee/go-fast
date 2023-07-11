## go-fast

FAST corner detector implementation with Golang

### Deploy

```shell script
git clone https://github.com/peterdee/go-fast
cd ./go-fast
mkdir samples
mkdir results
gvm use 1.20
```

### Launch

Modify `main.go` file

```go
// point radius (used for NMS clustering)
const RADIUS int = 15

// path to image file
const SAMPLE string = "samples/image.png"

// if result should be saved as grayscale
const SAVE_GRAYSCALE bool = false

// threshold for corner point determination
const THRESHOLD uint8 = 150
```

Run the code

```go
go run ./
```

Output file will be placed in `results` directory

### License

[MIT](LICENSE.md)
