import { isOfType } from 'typesafe-actions'
import { createSelector } from 'reselect'
import { takeEvery, put } from 'redux-saga/effects'

import { singular, plural } from 'pluralize'

import {
    reduce, find, omit, get, head, last, join, split, keys, size,
    snakeCase, toLower, includes, has
} from 'lodash'

import { grpc, RootState, AnyAction } from '../redux'

export enum Request {
    list    = 'LIST',
    get     = 'GET',
    create  = 'CREATE',
    update  = 'UPDATE',
    destroy = 'DESTROY'
}

export enum Status {
    requested = 'REQUESTED',
    succeeded = 'SUCCEEDED',
    failed    = 'FAILED'
}

export enum Resource {
    organization = 'organizations',
    project      = 'projects',
    site         = 'sites'
}

export type AnyResourceInstance = object

export type Selectors<T> = {
    getState: (state: RootState) => State<T> & any
    getAll: (state: RootState) => ResourcesList<T>
    countAll: (state: RootState) => number
    getByName: (name: string) => (state: RootState) => T | null
}

export type ActionTypes = {
    LIST_REQUESTED    : string,
    LIST_SUCCEEDED    : string,
    LIST_FAILED       : string,

    GET_REQUESTED     : string,
    GET_SUCCEEDED     : string,
    GET_FAILED        : string,

    CREATE_REQUESTED  : string,
    CREATE_SUCCEEDED  : string,
    CREATE_FAILED     : string,

    UPDATE_REQUESTED  : string,
    UPDATE_SUCCEEDED  : string,
    UPDATE_FAILED     : string,

    DESTROY_REQUESTED : string,
    DESTROY_SUCCEEDED : string,
    DESTROY_FAILED    : string
}

export type ResourcesList<T> = {
    [name: string]: T;
}

export type State<T> = {
    entries: ResourcesList<T>
}

export const initialState: State<AnyResourceInstance> = {
    entries: {}
}

export function createReducer(resource: Resource, actionTypes: ActionTypes) {
    return (state: State<AnyResourceInstance>, action: AnyAction) => {
        const response = get(action, 'payload')

        if (isOfType(actionTypes.LIST_SUCCEEDED, action)) {
            const entries = get(response, ['data', resource])
            return {
                ...state,
                entries: {
                    ...state.entries,
                    ...reduce(entries, (acc, entry) => ({
                        ...acc,
                        [entry.name]: {
                            ...get(acc, entry.name, {}),
                            ...entry
                        }
                    }), {})
                }
            }
        }

        if (isOfType([
            actionTypes.GET_SUCCEEDED,
            actionTypes.CREATE_SUCCEEDED,
            actionTypes.UPDATE_SUCCEEDED
        ], action)) {
            const entry = response.data
            return {
                ...state,
                entries: {
                    ...state.entries,
                    [entry.name]: {
                        ...get(state, ['entries', entry.name], {}),
                        ...entry
                    }
                }
            }
        }

        if (isOfType(actionTypes.DESTROY_SUCCEEDED, action)) {
            const entry = response.request.data
            return {
                ...state,
                entries: omit(state.entries, entry.name)
            }
        }

        return state
    }
}

export function createActionTypes(resource: Resource): ActionTypes {
    return reduce(Request, (accTypes, action) => {
        return {
            ...accTypes,
            ...reduce(Status, (acc, status) => {
                const key = createActionDescriptor(action, status)
                return {
                    ...acc,
                    [key]: `@ ${snakeCase(resource)} / ${key}`
                }
            }, {})
        }
    }, {} as any as ActionTypes)
}

export function createSelectors(resource: Resource): Selectors<AnyResourceInstance> {
    const getState = (state: RootState) => get(state, resource)
    const getAll = createSelector(
        getState,
        (state) => get(state, 'entries', {})
    )
    const getByName = (name: string) => createSelector(
        getAll,
        (entries) => get(entries, name, null)
    )
    const countAll = createSelector(
        getAll,
        (entries) => size(keys(entries))
    )
    return {
        getState, getAll, getByName, countAll
    }
}

export function* emitResourceActions(resource: Resource, actionTypes: ActionTypes) {
    yield takeEvery([
        grpc.INVOKED,
        grpc.SUCCEEDED,
        grpc.FAILED
    ], emitResourceAction, resource, actionTypes)
}

function* emitResourceAction(
    resource: Resource,
    actionTypes: ActionTypes,
    action: grpc.Actions
) {
    const method = isOfType(grpc.INVOKED, action)
        ? action.payload.method
        : action.payload.request.method

    const requestResource = getResourceFromMethod(method)
    const request = getRequestTypeFromMethod(method)
    const status = getStatusFromAction(action)

    if (requestResource !== resource) {
        return
    }

    if (!method || !request || !status) {
        return
    }

    const descriptor = createActionDescriptor(request, status)

    if (has(actionTypes, descriptor)) {
        yield put({ type: actionTypes[descriptor], payload: action.payload })
    }
}

function createActionDescriptor(request: Request, status: Status) {
    return join([request, status], '_')
}

function getRequestTypeFromMethod(methodName: string) {
    const requestType = toLower(head(split(snakeCase(methodName), '_')))

    if (!requestType) {
        return undefined
    }

    if (requestType === 'delete') {
        return Request.destroy
    }

    return get(Request, requestType, undefined)
}

function getResourceFromMethod(methodName: string) {
    const resourceName = toLower(last(split(snakeCase(methodName), '_')))

    if (!resourceName) {
        return undefined
    }

    return find(Resource, (resource) => (
        resource === plural(resourceName) ||
        resource === singular(resourceName)
    ))
}

function getStatusFromAction(action: grpc.Actions) {
    switch (action.type) {
        case grpc.INVOKED   : return Status.requested
        case grpc.SUCCEEDED : return Status.succeeded
        case grpc.FAILED    : return Status.failed
        default             : return undefined
    }
}
