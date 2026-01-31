package observability

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

func NewLogger() *zerolog.Logger {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", trimCallerPath(file, 3), line)
	}

	logger := zerolog.New(os.Stdout).With().CallerWithSkipFrameCount(2).Timestamp().Logger()

	return &logger
}

func trimCallerPath(path string, keepSegments int) string {
	if keepSegments <= 0 {
		return filepath.Base(path)
	}

	cutAt := lastNDirsIndex(path, keepSegments)
	if cutAt <= 0 {
		return path
	}

	return path[cutAt:]
}

func lastNDirsIndex(path string, keepSegments int) int {
	if keepSegments <= 0 {
		return 0
	}

	sep := byte(filepath.Separator)
	count := 0
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == sep {
			count++
			if count == keepSegments {
				return i + 1
			}
		}
	}

	return 0
}
