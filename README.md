# chanx

Unbounded chan.

[![License](https://img.shields.io/:license-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![GoDoc](https://godoc.org/github.com/smallnest/chanx?status.png)](http://godoc.org/github.com/smallnest/chanx)  [![travis](https://travis-ci.org/smallnest/chanx.svg?branch=master)](https://travis-ci.org/smallnest/chanx) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/chanx)](https://goreportcard.com/report/github.com/smallnest/chanx) [![coveralls](https://coveralls.io/repos/smallnest/chanx/badge.svg?branch=master&service=github)](https://coveralls.io/github/smallnest/chanx?branch=master) 

Refer to the below articles and issues:
1. https://github.com/golang/go/issues/20352
2. https://stackoverflow.com/questions/41906146/why-go-channels-limit-the-buffer-size
3. https://medium.com/capital-one-tech/building-an-unbounded-channel-in-go-789e175cd2cd
4. https://erikwinter.nl/articles/2020/channel-with-infinite-buffer-in-golang/


## Usage

```go
ch := NewUnboundedChan(1000)

go func() {
    ...
    ch.In <- ... // send values
    ...

    close(ch.In) // close In channel
}()


for v := range ch.Out { // read values
    fmt.Println(v)
}

```