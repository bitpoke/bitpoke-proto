import { ActionType, createAsyncAction, action as createAction } from 'typesafe-actions'
import { takeEvery, put, select, take, call } from 'redux-saga/effects'
import { SagaIterator, channel as createChannel } from 'redux-saga'
import { grpc } from 'grpc-web-client'
import { createSelector } from 'reselect'

import { reduce, get as _get, noop } from 'lodash'

import config from '../config'
import { watchChannel } from '../utils'
import { RootState, auth } from '../redux'

import {
    Organization,
    OrganizationsService,
    GetOrganizationRequest,
    CreateOrganizationRequest,
    UpdateOrganizationRequest,
    DeleteOrganizationRequest,
    ListOrganizationsRequest,
    ListOrganizationsResponse
} from '@presslabs/dashboard-proto'

const host: string = config.REACT_API_URL || 'http://localhost:8080'


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

export const GET_REQUESTED     = '@ organizations / GET_REQUESTED'
export const GET_SUCCEEDED     = '@ organizations / GET_SUCCEEDED'
export const LIST_REQUESTED    = '@ organizations / LIST_REQUESTED'
export const LIST_SUCCEEDED    = '@ organizations / LIST_SUCCEEDED'
export const CREATE_REQUESTED  = '@ organizations / CREATE_REQUESTED'
export const CREATE_SUCCEEDED  = '@ organizations / CREATE_SUCCEEDED'
export const UPDATE_REQUESTED  = '@ organizations / UPDATE_REQUESTED'
export const UPDATE_SUCCEEDED  = '@ organizations / UPDATE_SUCCEEDED'
export const DESTROY_REQUESTED = '@ organizations / DESTROY_REQUESTED'
export const DESTROY_SUCCEEDED = '@ organizations / DESTROY_SUCCEEDED'

export const get = () =>
    createAction(GET_REQUESTED)
export const getSuccess = (organization: Organization) =>
    createAction(GET_SUCCEEDED, organization)
export const list = () =>
    createAction(LIST_REQUESTED)
export const listSuccess = (response: ListOrganizationsResponse) =>
    createAction(LIST_SUCCEEDED, response)
export const create = (request: CreateOrganizationRequest) =>
    createAction(CREATE_REQUESTED, request)
export const createSuccess = (organization: Organization) =>
    createAction(CREATE_SUCCEEDED, organization)
export const update = (request: UpdateOrganizationRequest) =>
    createAction(UPDATE_REQUESTED, request)
export const updateSuccess = (organization: Organization) =>
    createAction(UPDATE_SUCCEEDED, organization)
export const destroy = (request: DeleteOrganizationRequest) =>
    createAction(DESTROY_REQUESTED, request)
export const destroySuccess = (organization: Organization) =>
    createAction(DESTROY_SUCCEEDED, organization)


const HANDLERS = {
    [GET_REQUESTED]: {
        service: OrganizationsService.GetOrganization,
        request: GetOrganizationRequest
    },
    [LIST_REQUESTED]: {
        service: OrganizationsService.ListOrganizations,
        request: ListOrganizationsRequest
    },
    [CREATE_REQUESTED]: {
        service: OrganizationsService.CreateOrganization,
        request: CreateOrganizationRequest
    },
    [UPDATE_REQUESTED]: {
        service: OrganizationsService.UpdateOrganization,
        request: UpdateOrganizationRequest
    },
    [DESTROY_REQUESTED]: {
        service: OrganizationsService.DeleteOrganization,
        request: DeleteOrganizationRequest
    }
}

const actions = {
    get, getSuccess,
    list, listSuccess,
    create, createSuccess,
    update, updateSuccess,
    destroy, destroySuccess
}


//
//  REDUCER

const initialState: State = {
    entries: {}
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case LIST_SUCCEEDED: {
            const response = action.payload
            return {
                ...state,
                entries: {
                    ...state.entries,
                    ...reduce(response.getOrganizationsList(), (acc, organization) => ({
                        ...acc,
                        [organization.getName()]: {
                            ..._get(acc, organization.getName(), {}),
                            ...organization.toObject()
                        }
                    }), {})
                }
            }
        }

        // case GET_SUCCEEDED:
        // case UPDATE_SUCCEEDED:
        case CREATE_SUCCEEDED: {
            const organization = action.payload
            return {
                ...state,
                entries: {
                    [organization.getName()]: {
                        ..._get(state, ['entries', organization.getName()], {}),
                        ...organization.toObject()
                    }
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
    yield takeEvery([LIST_REQUESTED, CREATE_REQUESTED], performRequest)
    yield watchChannel(channel)
}

function* performRequest(action: Actions) {
    const { service, request } = getHandlerForAction(action)
    const authorization = yield select(auth.getAuthorizationHeader)
    const metadata = { authorization }

    grpc.invoke(service, {
        host,
        metadata,
        request: new request(_get(action, 'payload')),
        onMessage: (response: ListOrganizationsResponse) =>
            channel.put(listSuccess(response)),
        onEnd: (response) => {
            console.log('>>>>>>>>>> onEnd response:')
            console.log(response)
        }
    })
}

function getHandlerForAction(action: Actions) {
    return _get(HANDLERS, action.type)
}

//
//  SELECTORS

export const getState = (state: RootState): State => state.organizations
export const getAll = createSelector(
    getState,
    (state: State): OrganizationsList => state.entries
)
