// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpslog

import (
	"context"
	"log/slog"
	"testing"

	"github.com/xmidt-org/wrp-go/v5"
)

// Test handlers are defined in handlers_test.go

func BenchmarkObserver_ObserveWRP(b *testing.B) {
	logger := slog.New(discardHandler{})

	msg := wrp.Message{
		Type:            wrp.SimpleRequestResponseMessageType,
		Source:          "dns:talaria.xmidt.example.com",
		Destination:     "event:device-status/mac:ffffffffdae4/online",
		TransactionUUID: "546514d4-9cb6-41c9-88ca-ccd4c130c525",
		ContentType:     "application/json",
		Payload:         []byte(`{"status":"online"}`),
	}

	tests := []struct {
		name   string
		fields []FieldOpt
	}{
		{
			name:   "single_field",
			fields: []FieldOpt{Source()},
		},
		{
			name:   "single_field_always",
			fields: []FieldOpt{SourceAlways()},
		},
		{
			name:   "three_fields",
			fields: []FieldOpt{Source(), Destination(), MessageType()},
		},
		{
			name:   "six_fields",
			fields: []FieldOpt{Source(), Destination(), MessageType(), TransactionUUID(), ContentType(), PayloadSize()},
		},
		{
			name: "all_fields",
			fields: []FieldOpt{
				MessageType(), Source(), Destination(), TransactionUUID(),
				ContentType(), Accept(), Status(), RequestDeliveryResponse(),
				Headers(), Metadata(), Path(), PayloadSize(), ServiceName(),
				URL(), PartnerIDs(), SessionID(), QualityOfService(),
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ob := Observer{
				Logger:  logger,
				Level:   slog.LevelInfo,
				Message: "wrp message",
				Fields:  tt.fields,
			}
			ctx := context.Background()

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				ob.ObserveWRP(ctx, msg)
			}
		})
	}
}

func BenchmarkObserver_ObserveWRP_Disabled(b *testing.B) {
	logger := slog.New(disabledHandler{})

	ob := Observer{
		Logger:  logger,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields:  []FieldOpt{Source(), Destination(), MessageType()},
	}

	msg := wrp.Message{
		Type:        wrp.SimpleRequestResponseMessageType,
		Source:      "dns:talaria.xmidt.example.com",
		Destination: "event:device-status/mac:ffffffffdae4/online",
	}

	ctx := context.Background()

	b.ReportAllocs()

	for b.Loop() {
		ob.ObserveWRP(ctx, msg)
	}
}

func BenchmarkObserver_ObserveWRP_NilLogger(b *testing.B) {
	ob := Observer{
		Logger:  nil,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields:  []FieldOpt{Source(), Destination(), MessageType()},
	}

	msg := wrp.Message{
		Type:        wrp.SimpleRequestResponseMessageType,
		Source:      "dns:talaria.xmidt.example.com",
		Destination: "event:device-status/mac:ffffffffdae4/online",
	}

	ctx := context.Background()

	b.ReportAllocs()

	for b.Loop() {
		ob.ObserveWRP(ctx, msg)
	}
}

func BenchmarkObserver_ObserveWRP_EmptyFields(b *testing.B) {
	logger := slog.New(discardHandler{})

	ob := Observer{
		Logger:  logger,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields: []FieldOpt{
			Source(), Destination(), MessageType(), Status(),
			Headers(), Metadata(), PartnerIDs(),
		},
	}

	// Message with mostly empty fields - tests the false return path
	msg := wrp.Message{
		Type:   wrp.SimpleRequestResponseMessageType,
		Source: "dns:talaria.xmidt.example.com",
	}

	ctx := context.Background()

	b.ReportAllocs()

	for b.Loop() {
		ob.ObserveWRP(ctx, msg)
	}
}
