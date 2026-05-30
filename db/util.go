package db

import (
	"runtime"
	"strings"
)

// callerExampleName walks the call stack looking for a frame whose file path
// sits inside an "examples/" directory and returns that directory's name
// (e.g. "01-context" from ".../examples/01-context/main.go").
// Returns an empty string if no such frame is found.
func callerExampleName() string {
	pcs := make([]uintptr, 16)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if idx := strings.Index(frame.File, "examples/"); idx != -1 {
			rest := frame.File[idx+len("examples/"):]
			if name, _, found := strings.Cut(rest, "/"); found {
				return name
			}
			return rest
		}
		if !more {
			break
		}
	}
	return ""
}
