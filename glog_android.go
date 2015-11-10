// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build android

package llog

// #cgo LDFLAGS: -llog
//
// #include <android/log.h>
import "C"

import (
	"time"
	"unsafe"
)

var ctag = C.CString("llog")

const maxLogSize = 1023 // from an off-hand comment in android/log.h

type androidLogger struct {
	prio C.int
}

func (l androidLogger) Flush() error { return nil }
func (l androidLogger) Sync() error  { return nil }
func (l androidLogger) Write(p []byte) (int, error) {
	n := len(p)
	for len(p) > 0 {
		p = l.writeOneLine(p)
	}
	return n, nil
}
func (l androidLogger) writeOneLine(p []byte) (remaining []byte) {
	n := len(p)
	var tmp byte
	if n > maxLogSize {
		n = maxLogSize
		remaining = p[n:]
		p = p[:n]
		tmp = remaining[0]
	}
	if p[n-1] != 0 {
		p = append(p, 0)
	}
	C.__android_log_write(l.prio, ctag, (*C.char)(unsafe.Pointer(&p[0])))
	if remaining != nil {
		// Restore the byte overwritten in append(p,0)
		remaining[0] = tmp
	}
	return remaining
}

func newFlushSyncWriter(l *Log, s Severity, now time.Time) (flushSyncWriter, error) {
	var prio C.int
	switch {
	case s <= InfoLog:
		prio = C.ANDROID_LOG_INFO
	case s <= WarningLog:
		prio = C.ANDROID_LOG_WARN
	case s <= ErrorLog:
		prio = C.ANDROID_LOG_ERROR
	case s >= FatalLog:
		prio = C.ANDROID_LOG_FATAL
	default:
		prio = C.ANDROID_LOG_DEFAULT
	}
	return androidLogger{prio}, nil
}
