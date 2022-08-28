package main

import "testing"

// this test would be detected as race condition
// go test -race .
func Test_updateMessage(t *testing.T) {
	msg = "hello, world!"

	wg.Add(2)
	go updateMessage("x")
	go updateMessage("Goodbye, cruel world!")
	wg.Wait()

	if msg != "Goodbye, cruel world!" {
		t.Error("incorrect value in msg")
	}
}
