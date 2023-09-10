# chanx

Unbounded chan with ringbuffer.

[![License](https://img.shields.io/:license-MIT-blue.svg)](https://opensource.org/licenses/MIT) ![GitHub](https://img.shields.io/github/license/smallnest/chanx) ![GitHub Action](https://github.com/smallnest/chanx/actions/workflows/go.yaml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/chanx)](https://goreportcard.com/report/github.com/smallnest/chanx)  [![GoDoc](https://godoc.org/github.com/smallnest/chanx?status.png)](http://godoc.org/github.com/smallnest/chanx) 

Refer to the below articles and issues:
1. https://github.com/golang/go/issues/20352
2. https://stackoverflow.com/questions/41906146/why-go-channels-limit-the-buffer-size
3. https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd
4. https://erikwinter.nl/articles/2020/channel-with-infinite-buffer-in-golang/


## Usage

If you want to use it with Go 1.17.x or below, you can use `github.com/smallnest/chanx@1.0.0`.
Since `github.com/smallnest/chanx@1.1.0`, it support Go generics.

```go
ch := NewUnboundedChan(1000)
// or ch := NewUnboundedChanSize(10,200,1000)

go func() {
    for ...... {
        ...
        ch.In <- ... // send values
        ...
    }

    close(ch.In) // close In channel
}()


for v := range ch.Out { // read values
    fmt.Println(v)
}

```