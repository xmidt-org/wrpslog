// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpslog

import (
	"context"
	"encoding/base64"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/wrp-go/v5"
)

// Test handlers are defined in handlers_test.go

func TestObserver_NilLogger(t *testing.T) {
	ob := Observer{
		Logger:  nil,
		Level:   slog.LevelInfo,
		Message: "test message",
		Fields:  []FieldOpt{Source(), Destination()},
	}

	// Should not panic
	ob.ObserveWRP(context.Background(), wrp.Message{})
}

func TestObserver_DisabledLevel(t *testing.T) {
	handler := newRecordHandler(slog.LevelWarn)
	logger := slog.New(handler)

	ob := Observer{
		Logger:  logger,
		Level:   slog.LevelInfo, // Below handler's level
		Message: "test message",
		Fields:  []FieldOpt{Source()},
	}

	ob.ObserveWRP(context.Background(), wrp.Message{Source: "test"})

	require.Empty(t, handler.records)
}

func TestObserver_FieldOptions(t *testing.T) {
	status := int64(200)
	rdr := int64(1)
	payload := []byte("test payload")

	fullMsg := wrp.Message{
		Type:                    wrp.SimpleRequestResponseMessageType,
		Source:                  "dns:test.example.com/service",
		Destination:             "mac:112233445566/config",
		TransactionUUID:         "uuid-1234",
		ContentType:             "application/json",
		Accept:                  "application/msgpack",
		Status:                  &status,
		RequestDeliveryResponse: &rdr,
		Headers:                 []string{"X-Header: value"},
		Metadata:                map[string]string{"key": "value"},
		Path:                    "/api/v1/test",
		Payload:                 payload,
		ServiceName:             "test-service",
		URL:                     "http://example.com",
		PartnerIDs:              []string{"partner1"},
		SessionID:               "session-abc",
		QualityOfService:        wrp.QOSMediumValue,
	}

	defaultFields := []FieldOpt{
		MessageType(),
		Source(),
		Destination(),
		TransactionUUID(),
		ContentType(),
		Accept(),
		Status(),
		RequestDeliveryResponse(),
		Headers(),
		Metadata(),
		Path(),
		PayloadAsBase64(),
		PayloadSize(),
		ServiceName(),
		URL(),
		PartnerIDs(),
		SessionID(),
		QualityOfService(),
	}

	alwaysFields := []FieldOpt{
		MessageType(),
		SourceAlways(),
		DestinationAlways(),
		TransactionUUIDAlways(),
		ContentTypeAlways(),
		AcceptAlways(),
		StatusAlways(),
		RequestDeliveryResponseAlways(),
		HeadersAlways(),
		MetadataAlways(),
		PathAlways(),
		PayloadAsBase64Always(),
		PayloadSizeAlways(),
		ServiceNameAlways(),
		URLAlways(),
		PartnerIDsAlways(),
		SessionIDAlways(),
		QualityOfServiceAlways(),
	}

	fullExpected := map[string]any{
		fMsgType:                 int64(wrp.SimpleRequestResponseMessageType),
		fSource:                  "dns:test.example.com/service",
		fDestination:             "mac:112233445566/config",
		fTransactionUUID:         "uuid-1234",
		fContentType:             "application/json",
		fAccept:                  "application/msgpack",
		fStatus:                  int64(200),
		fRequestDeliveryResponse: int64(1),
		fHeaders:                 []string{"X-Header: value"},
		fMetadata:                map[string]string{"key": "value"},
		fPath:                    "/api/v1/test",
		fPayload:                 base64.StdEncoding.EncodeToString(payload),
		fPayloadSize:             int64(len(payload)),
		fServiceName:             "test-service",
		fURL:                     "http://example.com",
		fPartnerIDs:              []string{"partner1"},
		fSessionID:               "session-abc",
		fQualityOfService:        int64(wrp.QOSMediumValue),
	}

	emptyExpected := map[string]any{
		fMsgType:                 int64(0),
		fSource:                  "",
		fDestination:             "",
		fTransactionUUID:         "",
		fContentType:             "",
		fAccept:                  "",
		fStatus:                  int64(0),
		fRequestDeliveryResponse: int64(0),
		fHeaders:                 []string(nil),
		fMetadata:                map[string]string(nil),
		fPath:                    "",
		fPayload:                 "",
		fPayloadSize:             int64(0),
		fServiceName:             "",
		fURL:                     "",
		fPartnerIDs:              []string(nil),
		fSessionID:               "",
		fQualityOfService:        int64(0),
	}

	tests := []struct {
		name     string
		fields   []FieldOpt
		msg      wrp.Message
		expected map[string]any
	}{
		{
			name:     "all_present_default",
			fields:   defaultFields,
			msg:      fullMsg,
			expected: fullExpected,
		},
		{
			name:     "all_present_always",
			fields:   alwaysFields,
			msg:      fullMsg,
			expected: fullExpected,
		},
		{
			name:     "all_empty_default",
			fields:   defaultFields,
			msg:      wrp.Message{},
			expected: map[string]any{fMsgType: int64(0), fQualityOfService: int64(0)}, // Only MessageType and QualityOfService logged
		},
		{
			name:     "all_empty_always",
			fields:   alwaysFields,
			msg:      wrp.Message{},
			expected: emptyExpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newRecordHandler(slog.LevelInfo)
			ob := Observer{
				Logger:  slog.New(handler),
				Level:   slog.LevelInfo,
				Message: "wrp message",
				Fields:  tt.fields,
			}

			ob.ObserveWRP(context.Background(), tt.msg)

			require.Len(t, handler.records, 1)
			attrs := handler.getAttrs(0)
			require.Len(t, attrs, len(tt.expected))
			for _, attr := range attrs {
				assert.Equal(t, tt.expected[attr.Key], attr.Value.Any(), "field %s", attr.Key)
			}
		})
	}
}

func TestObserver_MessageTypeAsString(t *testing.T) {
	handler := newRecordHandler(slog.LevelInfo)
	logger := slog.New(handler)

	ob := Observer{
		Logger:  logger,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields:  []FieldOpt{MessageTypeAsString()},
	}

	ob.ObserveWRP(context.Background(), wrp.Message{Type: wrp.SimpleRequestResponseMessageType})

	require.Len(t, handler.records, 1)
	attrs := handler.getAttrs(0)

	require.Len(t, attrs, 1)
	assert.Equal(t, fMsgType, attrs[0].Key)
	assert.Equal(t, wrp.SimpleRequestResponseMessageType.String(), attrs[0].Value.Any())
}

func TestObserver_DuplicateFieldsLastWins(t *testing.T) {
	handler := newRecordHandler(slog.LevelInfo)
	logger := slog.New(handler)

	ob := Observer{
		Logger:  logger,
		Level:   slog.LevelInfo,
		Message: "wrp message",
		Fields: []FieldOpt{
			MessageType(),         // First: numeric
			MessageTypeAsString(), // Second: string (should win)
		},
	}

	ob.ObserveWRP(context.Background(), wrp.Message{Type: wrp.SimpleRequestResponseMessageType})

	require.Len(t, handler.records, 1)
	attrs := handler.getAttrs(0)

	require.Len(t, attrs, 1)
	assert.Equal(t, fMsgType, attrs[0].Key)
	// String format should win since it was specified last
	assert.Equal(t, wrp.SimpleRequestResponseMessageType.String(), attrs[0].Value.Any())
}

var fieldMap = map[string]string{
	"Type":                    fMsgType,
	"Source":                  fSource,
	"Destination":             fDestination,
	"TransactionUUID":         fTransactionUUID,
	"ContentType":             fContentType,
	"Accept":                  fAccept,
	"Status":                  fStatus,
	"RequestDeliveryResponse": fRequestDeliveryResponse,
	"Headers":                 fHeaders,
	"Metadata":                fMetadata,
	"Path":                    fPath,
	"Payload":                 fPayload,
	"ServiceName":             fServiceName,
	"URL":                     fURL,
	"PartnerIDs":              fPartnerIDs,
	"SessionID":               fSessionID,
	"QualityOfService":        fQualityOfService,
}

func TestFieldOpt_JSONTags(t *testing.T) {
	// Spans and IncludeSpans are deprecated fields in wrp-go and intentionally
	// not supported by this package.
	ignored := map[string]struct{}{
		"Spans":        {},
		"IncludeSpans": {},
	}

	msgType := reflect.TypeOf(wrp.Message{})

	for fieldName, expectedTag := range fieldMap {
		t.Run(fieldName, func(t *testing.T) {
			field, found := msgType.FieldByName(fieldName)
			require.True(t, found, "Field '%s' not found in wrp.Message", fieldName)

			jsonTag := field.Tag.Get("json")
			list := strings.SplitN(jsonTag, ",", 2)
			require.NotEmpty(t, list[0], "Field '%s' does not have a JSON tag", fieldName)

			assert.Equal(t, list[0], expectedTag, "Constant for field '%s' does not match the JSON tag", fieldName)
		})
	}

	// Ensure all fields in wrp.Message are represented in the fieldMap
	for i := 0; i < msgType.NumField(); i++ {
		field := msgType.Field(i)
		if _, found := ignored[field.Name]; found {
			continue
		}

		assert.Contains(t, fieldMap, field.Name, "Field '%s' is not represented in the fieldMap", field.Name)
	}
}

func TestEnsureNoMissingFields(t *testing.T) {
	// Ensure that for every field in wrp.Message, there is a corresponding FieldOpt function
	msgType := reflect.TypeOf(wrp.Message{})

	for i := 0; i < msgType.NumField(); i++ {
		field := msgType.Field(i)

		if !field.IsExported() {
			continue
		}

		tag, found := fieldMap[field.Name]
		require.True(t, found, "Field '%s' is not represented in the fieldMap", field.Name)

		jsonTag := field.Tag.Get("json")
		list := strings.Split(jsonTag, ",")
		found = false
		for _, item := range list {
			if strings.TrimSpace(item) == tag {
				found = true
				break
			}
		}
		assert.True(t, found, "Constant for field '%s' (%s) does not match the JSON tag: %s", field.Name, tag, jsonTag)
	}
}
