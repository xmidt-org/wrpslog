// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

// Package wrpslog provides an observer that logs WRP messages using slog.
//
// The Observer implements wrp.Observer and logs configurable fields from WRP
// messages. Field names match the JSON tags used by wrp-go for consistency.
//
// # Usage
//
//	ob := &wrpslog.Observer{
//	    Logger:  slog.Default(),
//	    Level:   slog.LevelInfo,
//	    Message: "wrp message",
//	    Fields:  []wrpslog.FieldOpt{wrpslog.Source(), wrpslog.Destination()},
//	}
//
// # Field Selection
//
// Fields are configured using FieldOpt functions. Each field option sets a
// specific slot, so duplicate calls to the same field (e.g., MessageType and
// MessageTypeAsString) will overwrite - the last one wins.
//
// Empty or zero-value fields are automatically omitted from log output.
//
// # Performance
//
// The observer is designed for minimal allocations.
package wrpslog

import (
	"context"
	"log/slog"
	"sync"

	"github.com/xmidt-org/wrp-go/v5"
)

// Observer logs WRP messages using slog. It implements wrp.Observer.
//
// The observer must be used as a pointer (&Observer{}) to ensure proper
// initialization via sync.Once.
//
// Configuration fields (Logger, Level, Message, Fields) are read once on the
// first call to ObserveWRP. Modifications to these fields after the first call
// have no effect.
type Observer struct {
	// Logger is the slog.Logger to use. If nil, logging is skipped.
	Logger *slog.Logger

	// Level is the log level to use.
	Level slog.Level

	// Message is the log message text.
	Message string

	// Fields specifies which WRP fields to include in log output.
	// Each FieldOpt configures a specific field slot; duplicates overwrite.
	Fields []FieldOpt

	once   sync.Once
	fields [fieldCount]fieldFunc
}

var _ wrp.Observer = &Observer{}

func (ob *Observer) init() {
	ob.once.Do(func() {
		for _, opt := range ob.Fields {
			if opt != nil {
				opt(ob)
			}
		}
	})
}

// ObserveWRP logs information about the message being processed.
func (ob *Observer) ObserveWRP(ctx context.Context, msg wrp.Message) {
	if ob.Logger == nil {
		return
	}

	if !ob.Logger.Enabled(ctx, ob.Level) {
		return
	}

	ob.init()

	var idx int
	var attrs [fieldCount]slog.Attr
	for _, fn := range ob.fields {
		if fn != nil {
			if attr := fn(msg); attr.Key != "" {
				attrs[idx] = attr
				idx++
			}
		}
	}

	ob.Logger.LogAttrs(ctx, ob.Level, ob.Message, attrs[:idx]...)
}
