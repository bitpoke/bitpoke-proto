import { RPCImpl } from 'protobufjs'
import axios from 'axios'

import { parseInt, upperFirst, get } from 'lodash'

import { ChunkParser, ChunkType } from '../grpc/parser'
import config from '../../config'

const baseURL: string = config.REACT_API_URL || 'http://localhost:8080'
const timeout: number = config.REACT_API_DEFAULT_TIMEOUT || 10000

export enum Status {
  OK = 0,
  Canceled = 1,
  Unknown = 2,
  InvalidArgument = 3,
  DeadlineExceeded = 4,
  NotFound = 5,
  AlreadyExists = 6,
  PermissionDenied = 7,
  ResourceExhausted = 8,
  FailedPrecondition = 9,
  Aborted = 10,
  OutOfRange = 11,
  Unimplemented = 12,
  Internal = 13,
  Unavailable = 14,
  DataLoss = 15,
  Unauthenticated = 16
}

export function createTransport(serviceName: string): RPCImpl {
    const transport = axios.create({
        baseURL,
        timeout,
        headers : {
            'content-type' : 'application/grpc-web+proto',
            'x-grpc-web'   : '1'
        }
    })

    return async (method, requestData, callback) => {
        try {
            const response = await transport.request({
                url          : `${serviceName}/${upperFirst(method.name)}`,
                method       : 'POST',
                data         : frameRequest(requestData),
                responseType : 'arraybuffer'
            })

            const buffer = await response.data

            const status = parseInt(get(response.headers, 'grpc-status', 0), 0)
            const message = get(response.headers, 'grpc-message')

            if (response.status !== 200) {
                const error = new Error('Request failed')
                callback(error, null)
            }

            if (status !== Status.OK) {
                const error = new Error(message)
                callback(error, null)
            }
            else {
                const chunk = parseChunk(buffer)
                const data = new Uint8Array(get(chunk, 'data', []))
                callback(null, data)
            }
        }
        catch (error) {
            callback(error, null)
        }
    }
}

export function setMetadata(key: string, value: string) {
    axios.defaults.headers.common[key] = value
}

function parseChunk(buffer: ArrayBuffer) {
    if (buffer.byteLength === 0) {
        return null
    }
    return new ChunkParser()
        .parse(new Uint8Array(buffer))
        .find((chunk) => chunk.chunkType === ChunkType.MESSAGE)
}

function frameRequest(bytes: Uint8Array) {
    const frame = new ArrayBuffer(bytes.byteLength + 5)
    new DataView(frame, 1, 4).setUint32(0, bytes.length, false)
    new Uint8Array(frame, 5).set(bytes)
    return new Uint8Array(frame)
}
