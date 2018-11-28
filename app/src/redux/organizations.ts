import { ActionType, createAsyncAction, action as createAction } from 'typesafe-actions'
import { takeEvery, put, select, take, call } from 'redux-saga/effects'
import { SagaIterator, channel as createChannel } from 'redux-saga'
import { grpc } from 'grpc-web-client'
import { createSelector } from 'reselect'

import { reduce } from 'lodash'

import { watchChannel } from '../utils'
import { RootState, auth } from '../redux'

import { ListRequest } from '../proto/presslabs/dashboard/meta/v1/list_pb'
import { Organization } from '../proto/presslabs/dashboard/core/v1/organization_pb'
import { Organizations } from '../proto/presslabs/dashboard/core/v1/organization_pb_service'

const host: string = process.env.REACT_API_URL || 'http://localhost:9090'


//
//  TYPES

export type State = {
    readonly entries: OrganizationsList
}
export type Actions = ActionType<typeof actions>
export { Organization }

export type OrganizationsList = {
    [id: string]: Organization;
}

//
//  ACTIONS

export const LIST_REQUESTED = '@ organizations / LIST_REQUESTED'
export const LIST_SUCCEEDED = '@ organizations / LIST_SUCCEEDED'

export const RECEIVED       = '@ organizations / RECEIVED'

export const list = () => createAction(LIST_REQUESTED)
export const receive = (entry: Organization) => createAction(RECEIVED, entry)

const listRequest = {
    service: Organizations.List,
    request: new ListRequest()
}

const actions = {
    list,
    receive
}


//
//  REDUCER

const initialState: State = {
    entries: {}
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case RECEIVED: {
            const organization = action.payload
            return {
                ...state,
                entries: {
                    ...state.entries,
                    [organization.getId()]: organization
                }
            }
        }
    }

    return state
}


//
//  SAGA

const channel = createChannel()

export function* saga() {
    yield takeEvery(LIST_REQUESTED, performRequest)
    yield watchChannel(channel)
}

function* performRequest(action: ActionType<typeof list>) {
    const authorization = yield select(auth.getAuthorizationHeader)
    grpc.invoke(listRequest.service, {
        host,
        request: listRequest.request,
        metadata: { authorization },
        onMessage: (response: Organization) =>
            channel.put(receive(response)),
        onEnd: () =>
            console.log('ONEND!'),
        debug: true
    })
}


//
//  SELECTORS

export const getState = (state: RootState): State => state.organizations
export const getAll = createSelector(
    getState,
    (state: State): OrganizationsList => state.entries
)
