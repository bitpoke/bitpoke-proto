import { isOfType } from 'typesafe-actions'
import { createSelector } from 'reselect'
import { takeEvery, put, race, take, call } from 'redux-saga/effects'

import { singular, plural } from 'pluralize'

import {
    reduce, find, omit, get, head, last, join, split, keys, values as _values,
    size, replace, startCase, snakeCase, toLower, has, includes, startsWith, endsWith, isEmpty
} from 'lodash'

import { RootState, AnyAction, grpc, routing, forms, toasts } from '../redux'

import { Intent } from '@blueprintjs/core'

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

export type AnyResourceInstance = {
    name: string
}

export type Selectors<T> = {
    getState: (state: RootState) => State<T> & any
    getAll: (state: RootState) => ResourcesList<T>
    countAll: (state: RootState) => number
    getByName: (name: string) => (state: RootState) => T | null,
    getForURL: (url: routing.Path) => (state: RootState) => T | null
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

export type Actions = {
    list    : (payload: any) => AnyAction
    get     : (payload: any) => AnyAction
    create  : (payload: any) => AnyAction,
    update  : (payload: any) => AnyAction
    destroy : (payload: any) => AnyAction
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
        (state) => state.entries
    )
    const getByName = (name: string) => createSelector(
        getAll,
        (entries) => get(entries, name, null)
    )
    const countAll = createSelector(
        getAll,
        (entries) => size(keys(entries))
    )
    const getForURL = (fullURL: string) => createSelector(
        getAll,
        (entries) => find(entries, (entry) => startsWith(replace(fullURL, /^\//, ''), entry.name)) || null
    )
    return {
        getState, getAll, getByName, countAll, getForURL
    }
}

export function createFormHandler(
    formName: forms.Name,
    resource: Resource,
    actionTypes: ActionTypes,
    actions: Actions
) {
    return function* handleSubmit(action: forms.Actions) {
        const { resolve, reject, values } = action.payload
        const resourceName = singular(resource)
        const resourceDisplayName = startCase(resourceName)
        const entry = get(values, resourceName)

        if (isNewEntry(entry)) {
            yield put(actions.create(values))

            const { success } = yield race({
                success : take(actionTypes.CREATE_SUCCEEDED),
                failure : take(actionTypes.CREATE_FAILED)
            })

            if (success) {
                yield call(resolve)
                yield put(forms.reset(formName))
                toasts.show({ intent: Intent.SUCCESS, message: `${resourceDisplayName} created` })
            }
            else {
                yield call(reject, new forms.SubmissionError())
                toasts.show({ intent: Intent.DANGER, message: `${resourceDisplayName} create failed` })
            }
        }
        else {
            yield put(actions.update(values))

            const { success } = yield race({
                success : take(actionTypes.UPDATE_SUCCEEDED),
                failure : take(actionTypes.UPDATE_FAILED)
            })

            if (success) {
                yield call(resolve)
                yield put(forms.reset(formName))
                toasts.show({ intent: Intent.SUCCESS, message: `${resourceDisplayName} updated` })
            }
            else {
                yield call(reject, new forms.SubmissionError())
                toasts.show({ intent: Intent.DANGER, message: `${resourceDisplayName} update failed` })
            }
        }
    }
}

export function* emitResourceActions(resource: Resource, actionTypes: ActionTypes) {
    yield takeEvery([
        grpc.INVOKED,
        grpc.SUCCEEDED,
        grpc.FAILED
    ], emitResourceAction, resource, actionTypes)
}

export function isEmptyResponse({ type, payload }: { type: string, payload: grpc.Response }) {
    const { data, request } = payload

    if (!data || isEmpty(data)) {
        return true
    }

    const requestType = getRequestTypeFromMethodName(request.method)

    if (includes([Request.list, Request.get], requestType) && isEmpty(head(_values(data)))) {
        return true
    }

    return false
}

export function isNewEntry(entry: object | undefined) {
    return !has(entry, 'name')
}

function* emitResourceAction(
    resource: Resource,
    actionTypes: ActionTypes,
    action: grpc.Actions
) {
    const method = isOfType(grpc.INVOKED, action)
        ? action.payload.method
        : action.payload.request.method

    const requestResource = getResourceFromMethodName(method)
    const request = getRequestTypeFromMethodName(method)
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

function getRequestTypeFromMethodName(methodName: string): Request | null {
    const requestType = toLower(head(split(snakeCase(methodName), '_')))

    if (!requestType) {
        return null
    }

    if (requestType === 'delete') {
        return Request.destroy
    }

    return get(Request, requestType, null)
}

function getResourceFromMethodName(methodName: string): Resource | null {
    const resourceName = toLower(last(split(snakeCase(methodName), '_')))

    if (!resourceName) {
        return null
    }

    return find(Resource, (resource) => (
        resource === plural(resourceName) ||
        resource === singular(resourceName)
    )) || null
}

export function getStatusFromAction(action: AnyAction): Status | null {
    switch (action.type) {
        case grpc.INVOKED   : return Status.requested
        case grpc.SUCCEEDED : return Status.succeeded
        case grpc.FAILED    : return Status.failed
        default: {
            if (endsWith(action.type, Status.requested)) return Status.requested
            if (endsWith(action.type, Status.succeeded)) return Status.succeeded
            if (endsWith(action.type, Status.failed))    return Status.failed

            return null
        }
    }
}

export function getRequestTypeFromAction(action: AnyAction): Request | null {
    const method = get(action, 'payload.method', get(action, 'payload.request.method'))

    if (!method) {
        return null
    }

    return getRequestTypeFromMethodName(method)
}
