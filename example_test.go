// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpslog_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/xmidt-org/wrp-go/v5"
	"github.com/xmidt-org/wrpslog"
)

func Example() {
	// Create a logger that outputs to stdout for the example
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time and level for consistent example output
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// Create an observer that logs specific fields
	ob := &wrpslog.Observer{
		Logger:  logger,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields: []wrpslog.FieldOpt{
			wrpslog.MessageType(),
			wrpslog.Source(),
			wrpslog.Destination(),
			wrpslog.TransactionUUID(),
		},
	}

	// Create a WRP message
	msg := wrp.Message{
		Type:            wrp.SimpleRequestResponseMessageType,
		Source:          "dns:talaria.example.com/service",
		Destination:     "event:device-status/mac:112233445566/online",
		TransactionUUID: "546514d4-9cb6-41c9-88ca-ccd4c130c525",
	}

	// Log the message
	ob.ObserveWRP(context.Background(), msg)

	// Output:
	// msg="wrp message" msg_type=3 source=dns:talaria.example.com/service dest=event:device-status/mac:112233445566/online transaction_uuid=546514d4-9cb6-41c9-88ca-ccd4c130c525
}

func Example_alwaysFields() {
	// Create a logger that outputs to stdout for the example
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time and level for consistent example output
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// Use "Always" variants to log fields even when empty
	ob := &wrpslog.Observer{
		Logger:  logger,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields: []wrpslog.FieldOpt{
			wrpslog.MessageType(),
			wrpslog.Source(),
			wrpslog.DestinationAlways(), // Will log even if empty
			wrpslog.StatusAlways(),      // Will log even if nil
		},
	}

	// Message with some empty fields
	msg := wrp.Message{
		Type:   wrp.SimpleEventMessageType,
		Source: "dns:talaria.example.com/service",
		// Destination and Status are empty/nil
	}

	ob.ObserveWRP(context.Background(), msg)

	// Output:
	// msg="wrp message" msg_type=4 source=dns:talaria.example.com/service dest="" status=0
}
