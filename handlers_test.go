// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpslog

import (
	"context"
	"log/slog"
)

// discardHandler is a slog.Handler that discards all records without
// allocations. Always reports enabled.
type discardHandler struct{}

func (discardHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (d discardHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d discardHandler) WithGroup(string) slog.Handler           { return d }

// disabledHandler is a slog.Handler that always reports disabled. Useful for
// benchmarking the fast path when logging is turned off.
type disabledHandler struct{}

func (disabledHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (disabledHandler) Handle(context.Context, slog.Record) error { return nil }
func (d disabledHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d disabledHandler) WithGroup(string) slog.Handler           { return d }

// recordHandler is a slog.Handler that captures log records for test
// verification. Use newRecordHandler to create.
type recordHandler struct {
	records []slog.Record
	level   slog.Level
}

func newRecordHandler(level slog.Level) *recordHandler {
	return &recordHandler{
		records: make([]slog.Record, 0),
		level:   level,
	}
}

func (h *recordHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *recordHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}

func (h *recordHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *recordHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *recordHandler) getAttrs(index int) []slog.Attr {
	if index >= len(h.records) {
		return nil
	}
	var attrs []slog.Attr
	h.records[index].Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	return attrs
}
