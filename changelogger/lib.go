package changelogger

import (
	"github.com/sulaiman-coder/goeventbus"

	"github.com/khulnasoft-labs/changelogger/internal/bus"
	"github.com/khulnasoft-labs/changelogger/internal/log"
	"github.com/khulnasoft-labs/go-logger"
)

// SetLogger sets the logger object used for all logging calls.
func SetLogger(logger logger.Logger) {
	log.Log = logger
}

// SetBus sets the event bus for all published events onto (in-library subscriptions are not allowed).
func SetBus(b *eventbus.Bus) {
	bus.SetPublisher(b)
}
