// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpslog

import (
	"encoding/base64"
	"log/slog"

	"github.com/xmidt-org/wrp-go/v5"
)

// Field name constants matching WRP message JSON tags.
const (
	fMsgType                 = "msg_type"
	fSource                  = "source"
	fDestination             = "dest"
	fTransactionUUID         = "transaction_uuid"
	fContentType             = "content_type"
	fAccept                  = "accept"
	fStatus                  = "status"
	fRequestDeliveryResponse = "rdr"
	fHeaders                 = "headers"
	fMetadata                = "metadata"
	fPath                    = "path"
	fPayload                 = "payload"
	fPayloadSize             = "payload_size"
	fServiceName             = "service_name"
	fURL                     = "url"
	fPartnerIDs              = "partner_ids"
	fSessionID               = "session_id"
	fQualityOfService        = "qos"
)

// fieldIndex represents a specific field slot in the Fields array.
type fieldIndex int

const (
	idxMsgType fieldIndex = iota
	idxSource
	idxDestination
	idxTransactionUUID
	idxContentType
	idxAccept
	idxStatus
	idxRequestDeliveryResponse
	idxHeaders
	idxMetadata
	idxPath
	idxPayload
	idxPayloadSize
	idxServiceName
	idxURL
	idxPartnerIDs
	idxSessionID
	idxQualityOfService
	fieldCount // Total number of field slots
)

// FieldOpt configures a field to be logged by the Observer.
// Each FieldOpt sets a specific slot in the Observer's internal field array.
// Calling multiple options for the same field (e.g., MessageType and
// MessageTypeAsString) will overwrite - the last one wins.
type FieldOpt func(*Observer)

// fieldFunc extracts a field from a WRP message and returns it as an slog.Attr.
// Returns an slog.Attr with an empty Key to indicate the field should be skipped.
type fieldFunc func(wrp.Message) slog.Attr

// MessageType logs the message type as a number. This is an alias for MessageTypeAsNum.
func MessageType() FieldOpt {
	return MessageTypeAsNum()
}

// MessageTypeAsString logs the message type as a human-readable string.
// Uses the same slot as MessageType/MessageTypeAsNum.
func MessageTypeAsString() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxMsgType] = func(msg wrp.Message) slog.Attr {
			return slog.String(fMsgType, msg.Type.String())
		}
	}
}

// MessageTypeAsNum logs the message type as a number.
func MessageTypeAsNum() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxMsgType] = func(msg wrp.Message) slog.Attr {
			return slog.Int(fMsgType, int(msg.Type))
		}
	}
}

// Source logs the source of the message. Empty values are omitted.
func Source() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxSource] = func(msg wrp.Message) slog.Attr {
			if msg.Source == "" {
				return slog.Attr{}
			}
			return slog.String(fSource, msg.Source)
		}
	}
}

// SourceAlways logs the source of the message, even when empty.
func SourceAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxSource] = func(msg wrp.Message) slog.Attr {
			return slog.String(fSource, msg.Source)
		}
	}
}

// Destination logs the destination of the message. Empty values are omitted.
func Destination() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxDestination] = func(msg wrp.Message) slog.Attr {
			if msg.Destination == "" {
				return slog.Attr{}
			}
			return slog.String(fDestination, msg.Destination)
		}
	}
}

// DestinationAlways logs the destination of the message, even when empty.
func DestinationAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxDestination] = func(msg wrp.Message) slog.Attr {
			return slog.String(fDestination, msg.Destination)
		}
	}
}

// TransactionUUID logs the transaction UUID of the message. Empty values are omitted.
func TransactionUUID() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxTransactionUUID] = func(msg wrp.Message) slog.Attr {
			if msg.TransactionUUID == "" {
				return slog.Attr{}
			}
			return slog.String(fTransactionUUID, msg.TransactionUUID)
		}
	}
}

// TransactionUUIDAlways logs the transaction UUID of the message, even when empty.
func TransactionUUIDAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxTransactionUUID] = func(msg wrp.Message) slog.Attr {
			return slog.String(fTransactionUUID, msg.TransactionUUID)
		}
	}
}

// ContentType logs the content type of the message. Empty values are omitted.
func ContentType() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxContentType] = func(msg wrp.Message) slog.Attr {
			if msg.ContentType == "" {
				return slog.Attr{}
			}
			return slog.String(fContentType, msg.ContentType)
		}
	}
}

// ContentTypeAlways logs the content type of the message, even when empty.
func ContentTypeAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxContentType] = func(msg wrp.Message) slog.Attr {
			return slog.String(fContentType, msg.ContentType)
		}
	}
}

// Accept logs the accept header of the message. Empty values are omitted.
func Accept() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxAccept] = func(msg wrp.Message) slog.Attr {
			if msg.Accept == "" {
				return slog.Attr{}
			}
			return slog.String(fAccept, msg.Accept)
		}
	}
}

// AcceptAlways logs the accept header of the message, even when empty.
func AcceptAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxAccept] = func(msg wrp.Message) slog.Attr {
			return slog.String(fAccept, msg.Accept)
		}
	}
}

// Status logs the status of the message. Nil values are omitted.
func Status() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxStatus] = func(msg wrp.Message) slog.Attr {
			if msg.Status == nil {
				return slog.Attr{}
			}
			return slog.Int64(fStatus, *msg.Status)
		}
	}
}

// StatusAlways logs the status of the message, even when nil (logs 0).
func StatusAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxStatus] = func(msg wrp.Message) slog.Attr {
			if msg.Status == nil {
				return slog.Int64(fStatus, 0)
			}
			return slog.Int64(fStatus, *msg.Status)
		}
	}
}

// RequestDeliveryResponse logs the request delivery response of the message.
// Nil values are omitted.
func RequestDeliveryResponse() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxRequestDeliveryResponse] = func(msg wrp.Message) slog.Attr {
			if msg.RequestDeliveryResponse == nil {
				return slog.Attr{}
			}
			return slog.Int64(fRequestDeliveryResponse, *msg.RequestDeliveryResponse)
		}
	}
}

// RequestDeliveryResponseAlways logs the request delivery response of the
// message, even when nil (logs 0).
func RequestDeliveryResponseAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxRequestDeliveryResponse] = func(msg wrp.Message) slog.Attr {
			if msg.RequestDeliveryResponse == nil {
				return slog.Int64(fRequestDeliveryResponse, 0)
			}
			return slog.Int64(fRequestDeliveryResponse, *msg.RequestDeliveryResponse)
		}
	}
}

// Headers logs the headers of the message. Empty values are omitted.
func Headers() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxHeaders] = func(msg wrp.Message) slog.Attr {
			if len(msg.Headers) == 0 {
				return slog.Attr{}
			}
			return slog.Any(fHeaders, msg.Headers)
		}
	}
}

// HeadersAlways logs the headers of the message, even when empty.
func HeadersAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxHeaders] = func(msg wrp.Message) slog.Attr {
			return slog.Any(fHeaders, msg.Headers)
		}
	}
}

// Metadata logs the metadata of the message. Empty values are omitted.
func Metadata() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxMetadata] = func(msg wrp.Message) slog.Attr {
			if len(msg.Metadata) == 0 {
				return slog.Attr{}
			}
			return slog.Any(fMetadata, msg.Metadata)
		}
	}
}

// MetadataAlways logs the metadata of the message, even when empty.
func MetadataAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxMetadata] = func(msg wrp.Message) slog.Attr {
			return slog.Any(fMetadata, msg.Metadata)
		}
	}
}

// Path logs the path of the message. Empty values are omitted.
func Path() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPath] = func(msg wrp.Message) slog.Attr {
			if msg.Path == "" {
				return slog.Attr{}
			}
			return slog.String(fPath, msg.Path)
		}
	}
}

// PathAlways logs the path of the message, even when empty.
func PathAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPath] = func(msg wrp.Message) slog.Attr {
			return slog.String(fPath, msg.Path)
		}
	}
}

// PayloadAsBase64 logs the payload of the message as base64 encoded string.
// Empty values are omitted.
func PayloadAsBase64() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPayload] = func(msg wrp.Message) slog.Attr {
			if len(msg.Payload) == 0 {
				return slog.Attr{}
			}
			return slog.String(fPayload, base64.StdEncoding.EncodeToString(msg.Payload))
		}
	}
}

// PayloadAsBase64Always logs the payload of the message as base64 encoded
// string, even when empty.
func PayloadAsBase64Always() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPayload] = func(msg wrp.Message) slog.Attr {
			return slog.String(fPayload, base64.StdEncoding.EncodeToString(msg.Payload))
		}
	}
}

// PayloadSize logs the size of the payload of the message. Empty payloads are omitted.
func PayloadSize() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPayloadSize] = func(msg wrp.Message) slog.Attr {
			if len(msg.Payload) == 0 {
				return slog.Attr{}
			}
			return slog.Int(fPayloadSize, len(msg.Payload))
		}
	}
}

// PayloadSizeAlways logs the size of the payload of the message, even when empty.
func PayloadSizeAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPayloadSize] = func(msg wrp.Message) slog.Attr {
			return slog.Int(fPayloadSize, len(msg.Payload))
		}
	}
}

// ServiceName logs the service name of the message. Empty values are omitted.
func ServiceName() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxServiceName] = func(msg wrp.Message) slog.Attr {
			if msg.ServiceName == "" {
				return slog.Attr{}
			}
			return slog.String(fServiceName, msg.ServiceName)
		}
	}
}

// ServiceNameAlways logs the service name of the message, even when empty.
func ServiceNameAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxServiceName] = func(msg wrp.Message) slog.Attr {
			return slog.String(fServiceName, msg.ServiceName)
		}
	}
}

// URL logs the URL of the message. Empty values are omitted.
func URL() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxURL] = func(msg wrp.Message) slog.Attr {
			if msg.URL == "" {
				return slog.Attr{}
			}
			return slog.String(fURL, msg.URL)
		}
	}
}

// URLAlways logs the URL of the message, even when empty.
func URLAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxURL] = func(msg wrp.Message) slog.Attr {
			return slog.String(fURL, msg.URL)
		}
	}
}

// PartnerIDs logs the partner IDs of the message. Empty values are omitted.
func PartnerIDs() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPartnerIDs] = func(msg wrp.Message) slog.Attr {
			if len(msg.PartnerIDs) == 0 {
				return slog.Attr{}
			}
			return slog.Any(fPartnerIDs, msg.PartnerIDs)
		}
	}
}

// PartnerIDsAlways logs the partner IDs of the message, even when empty.
func PartnerIDsAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxPartnerIDs] = func(msg wrp.Message) slog.Attr {
			return slog.Any(fPartnerIDs, msg.PartnerIDs)
		}
	}
}

// SessionID logs the session ID of the message. Empty values are omitted.
func SessionID() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxSessionID] = func(msg wrp.Message) slog.Attr {
			if msg.SessionID == "" {
				return slog.Attr{}
			}
			return slog.String(fSessionID, msg.SessionID)
		}
	}
}

// SessionIDAlways logs the session ID of the message, even when empty.
func SessionIDAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxSessionID] = func(msg wrp.Message) slog.Attr {
			return slog.String(fSessionID, msg.SessionID)
		}
	}
}

// QualityOfService logs the quality of service of the message.
func QualityOfService() FieldOpt {
	return QualityOfServiceAlways()
}

// QualityOfServiceAlways logs the quality of service of the message, even when zero.
func QualityOfServiceAlways() FieldOpt {
	return func(ob *Observer) {
		ob.fields[idxQualityOfService] = func(msg wrp.Message) slog.Attr {
			return slog.Int(fQualityOfService, int(msg.QualityOfService))
		}
	}
}
