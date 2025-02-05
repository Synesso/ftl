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
// @generated from file opentelemetry/proto/resource/v1/resource.proto (package opentelemetry.proto.resource.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import { KeyValue } from "../../common/v1/common_pb.js";

/**
 * Resource information.
 *
 * @generated from message opentelemetry.proto.resource.v1.Resource
 */
export class Resource extends Message<Resource> {
  /**
   * Set of attributes that describe the resource.
   * Attribute keys MUST be unique (it is not allowed to have more than one
   * attribute with the same key).
   *
   * @generated from field: repeated opentelemetry.proto.common.v1.KeyValue attributes = 1;
   */
  attributes: KeyValue[] = [];

  /**
   * dropped_attributes_count is the number of dropped attributes. If the value is 0, then
   * no attributes were dropped.
   *
   * @generated from field: uint32 dropped_attributes_count = 2;
   */
  droppedAttributesCount = 0;

  constructor(data?: PartialMessage<Resource>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "opentelemetry.proto.resource.v1.Resource";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "attributes", kind: "message", T: KeyValue, repeated: true },
    { no: 2, name: "dropped_attributes_count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Resource {
    return new Resource().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Resource {
    return new Resource().fromJsonString(jsonString, options);
  }

  static equals(a: Resource | PlainMessage<Resource> | undefined, b: Resource | PlainMessage<Resource> | undefined): boolean {
    return proto3.util.equals(Resource, a, b);
  }
}

