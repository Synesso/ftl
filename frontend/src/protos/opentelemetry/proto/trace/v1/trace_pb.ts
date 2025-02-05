// Copyright 2019, OpenTelemetry Authors
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

// @generated by protoc-gen-es v1.7.0 with parameter "target=ts"
// @generated from file opentelemetry/proto/trace/v1/trace.proto (package opentelemetry.proto.trace.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64 } from "@bufbuild/protobuf";
import { Resource } from "../../resource/v1/resource_pb.js";
import { InstrumentationScope, KeyValue } from "../../common/v1/common_pb.js";

/**
 * TracesData represents the traces data that can be stored in a persistent storage,
 * OR can be embedded by other protocols that transfer OTLP traces data but do
 * not implement the OTLP protocol.
 *
 * The main difference between this message and collector protocol is that
 * in this message there will not be any "control" or "metadata" specific to
 * OTLP protocol.
 *
 * When new fields are added into this message, the OTLP request MUST be updated
 * as well.
 *
 * @generated from message opentelemetry.proto.trace.v1.TracesData
 */
export class TracesData extends Message<TracesData> {
  /**
   * An array of ResourceSpans.
   * For data coming from a single resource this array will typically contain
   * one element. Intermediary nodes that receive data from multiple origins
   * typically batch the data before forwarding further and in that case this
   * array will contain multiple elements.
   *
   * @generated from field: repeated opentelemetry.proto.trace.v1.ResourceSpans resource_spans = 1;
   */
  resourceSpans: ResourceSpans[] = [];

  constructor(data?: PartialMessage<TracesData>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.TracesData";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "resource_spans", kind: "message", T: ResourceSpans, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): TracesData {
    return new TracesData().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): TracesData {
    return new TracesData().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): TracesData {
    return new TracesData().fromJsonString(jsonString, options);
  }

  static equals(a: TracesData | PlainMessage<TracesData> | undefined, b: TracesData | PlainMessage<TracesData> | undefined): boolean {
    return proto3.util.equals(TracesData, a, b);
  }
}

/**
 * A collection of ScopeSpans from a Resource.
 *
 * @generated from message opentelemetry.proto.trace.v1.ResourceSpans
 */
export class ResourceSpans extends Message<ResourceSpans> {
  /**
   * The resource for the spans in this message.
   * If this field is not set then no resource info is known.
   *
   * @generated from field: opentelemetry.proto.resource.v1.Resource resource = 1;
   */
  resource?: Resource;

  /**
   * A list of ScopeSpans that originate from a resource.
   *
   * @generated from field: repeated opentelemetry.proto.trace.v1.ScopeSpans scope_spans = 2;
   */
  scopeSpans: ScopeSpans[] = [];

  /**
   * This schema_url applies to the data in the "resource" field. It does not apply
   * to the data in the "scope_spans" field which have their own schema_url field.
   *
   * @generated from field: string schema_url = 3;
   */
  schemaUrl = "";

  constructor(data?: PartialMessage<ResourceSpans>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.ResourceSpans";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "resource", kind: "message", T: Resource },
    { no: 2, name: "scope_spans", kind: "message", T: ScopeSpans, repeated: true },
    { no: 3, name: "schema_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ResourceSpans {
    return new ResourceSpans().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ResourceSpans {
    return new ResourceSpans().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ResourceSpans {
    return new ResourceSpans().fromJsonString(jsonString, options);
  }

  static equals(a: ResourceSpans | PlainMessage<ResourceSpans> | undefined, b: ResourceSpans | PlainMessage<ResourceSpans> | undefined): boolean {
    return proto3.util.equals(ResourceSpans, a, b);
  }
}

/**
 * A collection of Spans produced by an InstrumentationScope.
 *
 * @generated from message opentelemetry.proto.trace.v1.ScopeSpans
 */
export class ScopeSpans extends Message<ScopeSpans> {
  /**
   * The instrumentation scope information for the spans in this message.
   * Semantically when InstrumentationScope isn't set, it is equivalent with
   * an empty instrumentation scope name (unknown).
   *
   * @generated from field: opentelemetry.proto.common.v1.InstrumentationScope scope = 1;
   */
  scope?: InstrumentationScope;

  /**
   * A list of Spans that originate from an instrumentation scope.
   *
   * @generated from field: repeated opentelemetry.proto.trace.v1.Span spans = 2;
   */
  spans: Span[] = [];

  /**
   * This schema_url applies to all spans and span events in the "spans" field.
   *
   * @generated from field: string schema_url = 3;
   */
  schemaUrl = "";

  constructor(data?: PartialMessage<ScopeSpans>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.ScopeSpans";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "scope", kind: "message", T: InstrumentationScope },
    { no: 2, name: "spans", kind: "message", T: Span, repeated: true },
    { no: 3, name: "schema_url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ScopeSpans {
    return new ScopeSpans().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ScopeSpans {
    return new ScopeSpans().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ScopeSpans {
    return new ScopeSpans().fromJsonString(jsonString, options);
  }

  static equals(a: ScopeSpans | PlainMessage<ScopeSpans> | undefined, b: ScopeSpans | PlainMessage<ScopeSpans> | undefined): boolean {
    return proto3.util.equals(ScopeSpans, a, b);
  }
}

/**
 * A Span represents a single operation performed by a single component of the system.
 *
 * The next available field id is 17.
 *
 * @generated from message opentelemetry.proto.trace.v1.Span
 */
export class Span extends Message<Span> {
  /**
   * A unique identifier for a trace. All spans from the same trace share
   * the same `trace_id`. The ID is a 16-byte array. An ID with all zeroes OR
   * of length other than 16 bytes is considered invalid (empty string in OTLP/JSON
   * is zero-length and thus is also invalid).
   *
   * This field is required.
   *
   * @generated from field: bytes trace_id = 1;
   */
  traceId = new Uint8Array(0);

  /**
   * A unique identifier for a span within a trace, assigned when the span
   * is created. The ID is an 8-byte array. An ID with all zeroes OR of length
   * other than 8 bytes is considered invalid (empty string in OTLP/JSON
   * is zero-length and thus is also invalid).
   *
   * This field is required.
   *
   * @generated from field: bytes span_id = 2;
   */
  spanId = new Uint8Array(0);

  /**
   * trace_state conveys information about request position in multiple distributed tracing graphs.
   * It is a trace_state in w3c-trace-context format: https://www.w3.org/TR/trace-context/#tracestate-header
   * See also https://github.com/w3c/distributed-tracing for more details about this field.
   *
   * @generated from field: string trace_state = 3;
   */
  traceState = "";

  /**
   * The `span_id` of this span's parent span. If this is a root span, then this
   * field must be empty. The ID is an 8-byte array.
   *
   * @generated from field: bytes parent_span_id = 4;
   */
  parentSpanId = new Uint8Array(0);

  /**
   * A description of the span's operation.
   *
   * For example, the name can be a qualified method name or a file name
   * and a line number where the operation is called. A best practice is to use
   * the same display name at the same call point in an application.
   * This makes it easier to correlate spans in different traces.
   *
   * This field is semantically required to be set to non-empty string.
   * Empty value is equivalent to an unknown span name.
   *
   * This field is required.
   *
   * @generated from field: string name = 5;
   */
  name = "";

  /**
   * Distinguishes between spans generated in a particular context. For example,
   * two spans with the same name may be distinguished using `CLIENT` (caller)
   * and `SERVER` (callee) to identify queueing latency associated with the span.
   *
   * @generated from field: opentelemetry.proto.trace.v1.Span.SpanKind kind = 6;
   */
  kind = Span_SpanKind.UNSPECIFIED;

  /**
   * start_time_unix_nano is the start time of the span. On the client side, this is the time
   * kept by the local machine where the span execution starts. On the server side, this
   * is the time when the server's application handler starts running.
   * Value is UNIX Epoch time in nanoseconds since 00:00:00 UTC on 1 January 1970.
   *
   * This field is semantically required and it is expected that end_time >= start_time.
   *
   * @generated from field: fixed64 start_time_unix_nano = 7;
   */
  startTimeUnixNano = protoInt64.zero;

  /**
   * end_time_unix_nano is the end time of the span. On the client side, this is the time
   * kept by the local machine where the span execution ends. On the server side, this
   * is the time when the server application handler stops running.
   * Value is UNIX Epoch time in nanoseconds since 00:00:00 UTC on 1 January 1970.
   *
   * This field is semantically required and it is expected that end_time >= start_time.
   *
   * @generated from field: fixed64 end_time_unix_nano = 8;
   */
  endTimeUnixNano = protoInt64.zero;

  /**
   * attributes is a collection of key/value pairs. Note, global attributes
   * like server name can be set using the resource API. Examples of attributes:
   *
   *     "/http/user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
   *     "/http/server_latency": 300
   *     "example.com/myattribute": true
   *     "example.com/score": 10.239
   *
   * The OpenTelemetry API specification further restricts the allowed value types:
   * https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/common/README.md#attribute
   * Attribute keys MUST be unique (it is not allowed to have more than one
   * attribute with the same key).
   *
   * @generated from field: repeated opentelemetry.proto.common.v1.KeyValue attributes = 9;
   */
  attributes: KeyValue[] = [];

  /**
   * dropped_attributes_count is the number of attributes that were discarded. Attributes
   * can be discarded because their keys are too long or because there are too many
   * attributes. If this value is 0, then no attributes were dropped.
   *
   * @generated from field: uint32 dropped_attributes_count = 10;
   */
  droppedAttributesCount = 0;

  /**
   * events is a collection of Event items.
   *
   * @generated from field: repeated opentelemetry.proto.trace.v1.Span.Event events = 11;
   */
  events: Span_Event[] = [];

  /**
   * dropped_events_count is the number of dropped events. If the value is 0, then no
   * events were dropped.
   *
   * @generated from field: uint32 dropped_events_count = 12;
   */
  droppedEventsCount = 0;

  /**
   * links is a collection of Links, which are references from this span to a span
   * in the same or different trace.
   *
   * @generated from field: repeated opentelemetry.proto.trace.v1.Span.Link links = 13;
   */
  links: Span_Link[] = [];

  /**
   * dropped_links_count is the number of dropped links after the maximum size was
   * enforced. If this value is 0, then no links were dropped.
   *
   * @generated from field: uint32 dropped_links_count = 14;
   */
  droppedLinksCount = 0;

  /**
   * An optional final status for this span. Semantically when Status isn't set, it means
   * span's status code is unset, i.e. assume STATUS_CODE_UNSET (code = 0).
   *
   * @generated from field: opentelemetry.proto.trace.v1.Status status = 15;
   */
  status?: Status;

  constructor(data?: PartialMessage<Span>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.Span";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "trace_id", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "span_id", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 3, name: "trace_state", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "parent_span_id", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 5, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 6, name: "kind", kind: "enum", T: proto3.getEnumType(Span_SpanKind) },
    { no: 7, name: "start_time_unix_nano", kind: "scalar", T: 6 /* ScalarType.FIXED64 */ },
    { no: 8, name: "end_time_unix_nano", kind: "scalar", T: 6 /* ScalarType.FIXED64 */ },
    { no: 9, name: "attributes", kind: "message", T: KeyValue, repeated: true },
    { no: 10, name: "dropped_attributes_count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 11, name: "events", kind: "message", T: Span_Event, repeated: true },
    { no: 12, name: "dropped_events_count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 13, name: "links", kind: "message", T: Span_Link, repeated: true },
    { no: 14, name: "dropped_links_count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 15, name: "status", kind: "message", T: Status },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Span {
    return new Span().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Span {
    return new Span().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Span {
    return new Span().fromJsonString(jsonString, options);
  }

  static equals(a: Span | PlainMessage<Span> | undefined, b: Span | PlainMessage<Span> | undefined): boolean {
    return proto3.util.equals(Span, a, b);
  }
}

/**
 * SpanKind is the type of span. Can be used to specify additional relationships between spans
 * in addition to a parent/child relationship.
 *
 * @generated from enum opentelemetry.proto.trace.v1.Span.SpanKind
 */
export enum Span_SpanKind {
  /**
   * Unspecified. Do NOT use as default.
   * Implementations MAY assume SpanKind to be INTERNAL when receiving UNSPECIFIED.
   *
   * @generated from enum value: SPAN_KIND_UNSPECIFIED = 0;
   */
  UNSPECIFIED = 0,

  /**
   * Indicates that the span represents an internal operation within an application,
   * as opposed to an operation happening at the boundaries. Default value.
   *
   * @generated from enum value: SPAN_KIND_INTERNAL = 1;
   */
  INTERNAL = 1,

  /**
   * Indicates that the span covers server-side handling of an RPC or other
   * remote network request.
   *
   * @generated from enum value: SPAN_KIND_SERVER = 2;
   */
  SERVER = 2,

  /**
   * Indicates that the span describes a request to some remote service.
   *
   * @generated from enum value: SPAN_KIND_CLIENT = 3;
   */
  CLIENT = 3,

  /**
   * Indicates that the span describes a producer sending a message to a broker.
   * Unlike CLIENT and SERVER, there is often no direct critical path latency relationship
   * between producer and consumer spans. A PRODUCER span ends when the message was accepted
   * by the broker while the logical processing of the message might span a much longer time.
   *
   * @generated from enum value: SPAN_KIND_PRODUCER = 4;
   */
  PRODUCER = 4,

  /**
   * Indicates that the span describes consumer receiving a message from a broker.
   * Like the PRODUCER kind, there is often no direct critical path latency relationship
   * between producer and consumer spans.
   *
   * @generated from enum value: SPAN_KIND_CONSUMER = 5;
   */
  CONSUMER = 5,
}
// Retrieve enum metadata with: proto3.getEnumType(Span_SpanKind)
proto3.util.setEnumType(Span_SpanKind, "opentelemetry.proto.trace.v1.Span.SpanKind", [
  { no: 0, name: "SPAN_KIND_UNSPECIFIED" },
  { no: 1, name: "SPAN_KIND_INTERNAL" },
  { no: 2, name: "SPAN_KIND_SERVER" },
  { no: 3, name: "SPAN_KIND_CLIENT" },
  { no: 4, name: "SPAN_KIND_PRODUCER" },
  { no: 5, name: "SPAN_KIND_CONSUMER" },
]);

/**
 * Event is a time-stamped annotation of the span, consisting of user-supplied
 * text description and key-value pairs.
 *
 * @generated from message opentelemetry.proto.trace.v1.Span.Event
 */
export class Span_Event extends Message<Span_Event> {
  /**
   * time_unix_nano is the time the event occurred.
   *
   * @generated from field: fixed64 time_unix_nano = 1;
   */
  timeUnixNano = protoInt64.zero;

  /**
   * name of the event.
   * This field is semantically required to be set to non-empty string.
   *
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * attributes is a collection of attribute key/value pairs on the event.
   * Attribute keys MUST be unique (it is not allowed to have more than one
   * attribute with the same key).
   *
   * @generated from field: repeated opentelemetry.proto.common.v1.KeyValue attributes = 3;
   */
  attributes: KeyValue[] = [];

  /**
   * dropped_attributes_count is the number of dropped attributes. If the value is 0,
   * then no attributes were dropped.
   *
   * @generated from field: uint32 dropped_attributes_count = 4;
   */
  droppedAttributesCount = 0;

  constructor(data?: PartialMessage<Span_Event>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.Span.Event";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "time_unix_nano", kind: "scalar", T: 6 /* ScalarType.FIXED64 */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "attributes", kind: "message", T: KeyValue, repeated: true },
    { no: 4, name: "dropped_attributes_count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Span_Event {
    return new Span_Event().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Span_Event {
    return new Span_Event().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Span_Event {
    return new Span_Event().fromJsonString(jsonString, options);
  }

  static equals(a: Span_Event | PlainMessage<Span_Event> | undefined, b: Span_Event | PlainMessage<Span_Event> | undefined): boolean {
    return proto3.util.equals(Span_Event, a, b);
  }
}

/**
 * A pointer from the current span to another span in the same trace or in a
 * different trace. For example, this can be used in batching operations,
 * where a single batch handler processes multiple requests from different
 * traces or when the handler receives a request from a different project.
 *
 * @generated from message opentelemetry.proto.trace.v1.Span.Link
 */
export class Span_Link extends Message<Span_Link> {
  /**
   * A unique identifier of a trace that this linked span is part of. The ID is a
   * 16-byte array.
   *
   * @generated from field: bytes trace_id = 1;
   */
  traceId = new Uint8Array(0);

  /**
   * A unique identifier for the linked span. The ID is an 8-byte array.
   *
   * @generated from field: bytes span_id = 2;
   */
  spanId = new Uint8Array(0);

  /**
   * The trace_state associated with the link.
   *
   * @generated from field: string trace_state = 3;
   */
  traceState = "";

  /**
   * attributes is a collection of attribute key/value pairs on the link.
   * Attribute keys MUST be unique (it is not allowed to have more than one
   * attribute with the same key).
   *
   * @generated from field: repeated opentelemetry.proto.common.v1.KeyValue attributes = 4;
   */
  attributes: KeyValue[] = [];

  /**
   * dropped_attributes_count is the number of dropped attributes. If the value is 0,
   * then no attributes were dropped.
   *
   * @generated from field: uint32 dropped_attributes_count = 5;
   */
  droppedAttributesCount = 0;

  constructor(data?: PartialMessage<Span_Link>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.Span.Link";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "trace_id", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 2, name: "span_id", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 3, name: "trace_state", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "attributes", kind: "message", T: KeyValue, repeated: true },
    { no: 5, name: "dropped_attributes_count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Span_Link {
    return new Span_Link().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Span_Link {
    return new Span_Link().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Span_Link {
    return new Span_Link().fromJsonString(jsonString, options);
  }

  static equals(a: Span_Link | PlainMessage<Span_Link> | undefined, b: Span_Link | PlainMessage<Span_Link> | undefined): boolean {
    return proto3.util.equals(Span_Link, a, b);
  }
}

/**
 * The Status type defines a logical error model that is suitable for different
 * programming environments, including REST APIs and RPC APIs.
 *
 * @generated from message opentelemetry.proto.trace.v1.Status
 */
export class Status extends Message<Status> {
  /**
   * A developer-facing human readable error message.
   *
   * @generated from field: string message = 2;
   */
  message = "";

  /**
   * The status code.
   *
   * @generated from field: opentelemetry.proto.trace.v1.Status.StatusCode code = 3;
   */
  code = Status_StatusCode.UNSET;

  constructor(data?: PartialMessage<Status>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.trace.v1.Status";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 2, name: "message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "code", kind: "enum", T: proto3.getEnumType(Status_StatusCode) },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Status {
    return new Status().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Status {
    return new Status().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Status {
    return new Status().fromJsonString(jsonString, options);
  }

  static equals(a: Status | PlainMessage<Status> | undefined, b: Status | PlainMessage<Status> | undefined): boolean {
    return proto3.util.equals(Status, a, b);
  }
}

/**
 * For the semantics of status codes see
 * https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md#set-status
 *
 * @generated from enum opentelemetry.proto.trace.v1.Status.StatusCode
 */
export enum Status_StatusCode {
  /**
   * The default status.
   *
   * @generated from enum value: STATUS_CODE_UNSET = 0;
   */
  UNSET = 0,

  /**
   * The Span has been validated by an Application developer or Operator to
   * have completed successfully.
   *
   * @generated from enum value: STATUS_CODE_OK = 1;
   */
  OK = 1,

  /**
   * The Span contains an error.
   *
   * @generated from enum value: STATUS_CODE_ERROR = 2;
   */
  ERROR = 2,
}
// Retrieve enum metadata with: proto3.getEnumType(Status_StatusCode)
proto3.util.setEnumType(Status_StatusCode, "opentelemetry.proto.trace.v1.Status.StatusCode", [
  { no: 0, name: "STATUS_CODE_UNSET" },
  { no: 1, name: "STATUS_CODE_OK" },
  { no: 2, name: "STATUS_CODE_ERROR" },
]);

