import { channel as createChannel } from 'redux-saga'
import { takeEvery, call, put } from 'redux-saga/effects'
import { ActionType, action as createAction } from 'typesafe-actions'
import { createSelector } from 'reselect'

import { get as _get, findIndex, concat, isEqual, isEmpty } from 'lodash'

import { RootState } from '../redux'

import { watchChannel } from '../utils'
import { createTransport, setMetadataHeader } from '../utils/grpc/transport'

export { createTransport }

export type State = {
    ongoingRequests: [],
    metadata: {
        [key: string]: string
    }
}

export type Actions = ActionType<typeof actions>

export type Request = {
    method  : string,
    data    : object,
    service : object
}

export type Response = {
    data?   : object | null,
    error?  : Error,
    request : Request
}

export type Metadata = {
    key   : string,
    value : string
}

//
//  ACTIONS

export const INVOKED      = '@ grpc / INVOKED'
export const SUCCEEDED    = '@ grpc / SUCCEEDED'
export const FAILED       = '@ grpc / FAILED'
export const METADATA_SET = '@ grpc / METADATA_SET'

export const invoke = (payload: Request) => createAction(INVOKED, payload)
export const success = (payload: Response) => createAction(SUCCEEDED, payload)
export const fail = (payload: Response) => createAction(FAILED, payload)
export const setMetadata = (payload: Metadata) => createAction(METADATA_SET, payload)

const actions = {
    invoke,
    success,
    fail,
    setMetadata
}

//
//  SAGA

const channel = createChannel()

export function* saga() {
    yield takeEvery(INVOKED, performRequest)
    yield takeEvery(METADATA_SET, updateTransportMetadata)
    yield watchChannel(channel)
}

function* performRequest(action: ActionType<typeof invoke>) {
    const request = action.payload
    const { service, method, data } = request

    try {
        service[method](data, (error: Error|null, responseData: object|null) => {
            if (error) {
                channel.put(fail({ request, error }))
                return false
            }

            channel.put(success({ request, data: responseData }))
            return true
        })
    }
    catch (error) {
        channel.put(fail({ request, error }))
        return false
    }
}

function* updateTransportMetadata(action: ActionType<typeof setMetadata>) {
    const metadata = action.payload
    const { key, value } = metadata
    yield call(setMetadataHeader, key, value)
}


//
//  REDUCER

const initialState: State = {
    ongoingRequests: [],
    metadata: {}
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case INVOKED: {
            const request = action.payload
            return {
                ...state,
                ongoingRequests: concat(state.ongoingRequests, [request])
            }
        }

        case SUCCEEDED:
        case FAILED: {
            const request = action.payload.request
            const index = findIndex(state.ongoingRequests, (req) => isEqual(req, request))

            if (index < 0) {
                return state
            }

            const ongoingRequests = [...state.ongoingRequests]
            ongoingRequests.splice(index, 1)

            return {
                ...state,
                ongoingRequests
            }
        }

        case METADATA_SET: {
            const { key, value } = action.payload
            return {
                ...state,
                metadata: {
                    [key]: value
                }
            }
        }

        default:
            return state
    }
}

//
//  SELECTORS

export const getState = (state: RootState) => state.grpc
export const isLoading = createSelector(
    getState,
    (state) => !isEmpty(state.ongoingRequests)
)
export const getMetadata = (key: string) => createSelector(
    getState,
    (state) => _get(state, ['metadata', key])
)
