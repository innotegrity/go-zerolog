//go:build !binary_log
// +build !binary_log

package pkgerrors

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"go.innotegrity.dev/zerolog"
)

func TestLogStack(t *testing.T) {
	zerolog.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out)

	err := fmt.Errorf("from error: %w", errors.New("error message"))
	log.Log().Stack().Err(err).Msg("")

	got := out.String()
	//want := `\{"stack":\[\{"func":"TestLogStack","line":"21","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	want := `\{"stack":\[\{"func":"TestLogStack","line":"22","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	/** BEGIN CUSTOM CODE **/
	/** END CUSTOM CODE **/
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLogStackFromContext(t *testing.T) {
	zerolog.ErrorStackMarshaler = MarshalStack

	out := &bytes.Buffer{}
	log := zerolog.New(out).With().Stack().Logger() // calling Stack() on log context instead of event

	err := fmt.Errorf("from error: %w", errors.New("error message"))
	log.Log().Err(err).Msg("") // not explicitly calling Stack()

	got := out.String()
	/** BEGIN CUSTOM CODE **/
	//want := `\{"stack":\[\{"func":"TestLogStackFromContext","line":"37","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	want := `\{"stack":\[\{"func":"TestLogStackFromContext","line":"41","source":"stacktrace_test.go"\},.*\],"error":"from error: error message"\}\n`
	/** END CUSTOM CODE **/
	if ok, _ := regexp.MatchString(want, got); !ok {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func BenchmarkLogStack(b *testing.B) {
	zerolog.ErrorStackMarshaler = MarshalStack
	out := &bytes.Buffer{}
	log := zerolog.New(out)
	err := errors.Wrap(errors.New("error message"), "from error")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		log.Log().Stack().Err(err).Msg("")
		out.Reset()
	}
}
