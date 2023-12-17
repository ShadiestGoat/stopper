# Stopper

[![Go Reference](https://pkg.go.dev/badge/github.com/shadiestgoat/stopper.svg)](https://pkg.go.dev/github.com/shadiestgoat/stopper)

An easy to use modules to stop your go routines.

This module targets a very specific yet very common issue - you have a lot of go routines, and you need to stop them all as part of your app's graceful shutdown.

## Install

```bash
go get github.com/shadiestgoat/stopper
```

## How to use

In this module, there are 2 main structures - a sender and a receiver. The idea is that in your initialization, you create a sender using `NewSender`, then register a stopper per go routine you need to stop. The receiver will receive a stop signal over the `Receiver.C` channel, which will signal your go routine to stop what its doing.

Realistically there should be 1 sender, but if you need multiple you can create more.

## Example

```go
package main

func Ticker(s *stopper.Receiver, v int, name string) {
    t := time.NewTicker(v * time.Second)
    defer s.Done()

    for {
        select {
        case <- t.C:
            fmt.Printf("Aha! %v rules!\n", name)
        case <- s.C:
            fmt.Printf("Oh no! %v closing...\n", name)
            doLongOperationLikeWritingToAFileHere()
            return
        }
    }
}

func main() {
    s := stopper.NewAsync(nil)

    go Ticker(s.Register("foo"), 1, "foo")
    go Ticker(s.Register("bar"), 2, "bar")
    go Ticker(s.Register("abc"), 3, "abc")

	endApp := make(chan os.Signal, 2)

	signal.Notify(endApp, os.Interrupt)

    log("Closing all my tickers...")
    s.Stop()
    log("All my tickers are now closed, and I can rest easily since all my files are properly saved!")
}
```