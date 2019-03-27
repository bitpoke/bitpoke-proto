import { SagaIterator, channel as createChannel } from 'redux-saga'
import { takeEvery, put, select, take, call } from 'redux-saga/effects'
import { ActionType, action as createAction } from 'typesafe-actions'
import { createSelector } from 'reselect'

import {
    reduce, get as _get, findIndex, compact, concat, join, noop, size,
    snakeCase, startCase, toUpper, toLower, includes, isEqual, isString, isEmpty
} from 'lodash'

import { RootState } from '../redux'

import { watchChannel } from '../utils'
export { createTransport, setMetadata } from '../utils/grpc/transport'

export type State = {
    ongoingRequests: []
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

//
//  ACTIONS

export const INVOKED   = '@ grpc / INVOKED'
export const SUCCEEDED = '@ grpc / SUCCEEDED'
export const FAILED    = '@ grpc / FAILED'

export const invoke = (payload: Request) => createAction(INVOKED, payload)
export const success = (payload: Response) => createAction(SUCCEEDED, payload)
export const fail = (payload: Response) => createAction(FAILED, payload)

const types = {
    INVOKED,
    SUCCEEDED,
    FAILED
}

const actions = {
    invoke,
    success,
    fail
}

//
//  SAGA

const channel = createChannel()

export function* saga() {
    yield takeEvery(INVOKED, performRequest)
    yield watchChannel(channel)
}

function* performRequest(action: ReturnType<typeof invoke>) {
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


//
//  REDUCER

const initialState: State = {
    ongoingRequests: []
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
