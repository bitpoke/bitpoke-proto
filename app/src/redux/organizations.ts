import { ActionType, createAsyncAction, action as createAction, isOfType } from 'typesafe-actions'
import {
    takeEvery, takeLatest, fork, put, take, call, race, select as _select
} from 'redux-saga/effects'
import { SagaIterator, channel as createChannel } from 'redux-saga'
import { createSelector } from 'reselect'
import {
    reduce, find, head, values, join, noop,
    snakeCase, toUpper, has, includes, get as _get, isEmpty, isEqual
} from 'lodash'

import URI from 'urijs'

import { RootState, app, api, grpc, routing } from '../redux'

import { presslabs } from '@presslabs/dashboard-proto'

const {
    Organization,
    OrganizationsService,
    ListOrganizationsRequest,
    ListOrganizationsResponse,
    GetOrganizationRequest,
    CreateOrganizationRequest,
    UpdateOrganizationRequest,
    DeleteOrganizationRequest
} = presslabs.dashboard.organizations.v1

export {
    Organization,
    OrganizationsService,
    ListOrganizationsRequest,
    ListOrganizationsResponse,
    GetOrganizationRequest,
    CreateOrganizationRequest,
    UpdateOrganizationRequest,
    DeleteOrganizationRequest
}


//
//  TYPES

export type Organization =
    presslabs.dashboard.organizations.v1.Organization

export type IOrganization =
    presslabs.dashboard.organizations.v1.IOrganization

export type ListOrganizationsResponse =
    presslabs.dashboard.organizations.v1.ListOrganizationsRequest

export type IListOrganizationsRequest =
    presslabs.dashboard.organizations.v1.IListOrganizationsRequest

export type IGetOrganizationRequest =
    presslabs.dashboard.organizations.v1.IGetOrganizationRequest

export type ICreateOrganizationRequest =
    presslabs.dashboard.organizations.v1.ICreateOrganizationRequest

export type IUpdateOrganizationRequest =
    presslabs.dashboard.organizations.v1.IUpdateOrganizationRequest

export type IDeleteOrganizationRequest =
    presslabs.dashboard.organizations.v1.IDeleteOrganizationRequest

export type IListOrganizationsResponse =
    presslabs.dashboard.organizations.v1.IListOrganizationsResponse

export type State = {
    current: string | null
} & api.State<IOrganization>

export type Actions = ActionType<typeof actions>

const resource = api.Resource.organization

const service = OrganizationsService.create(
    grpc.createTransport('presslabs.dashboard.organizations.v1.OrganizationsService')
)


//
//  ACTIONS

export const SELECTED = '@ organizations / SELECTED'

export const select = (payload: IOrganization) => createAction(SELECTED, payload)

export const get = (payload: IGetOrganizationRequest) => grpc.invoke({
    service,
    method : 'getOrganization',
    data   : GetOrganizationRequest.create(payload)
})

export const list = (payload?: IListOrganizationsRequest) => grpc.invoke({
    service,
    method : 'listOrganizations',
    data   : ListOrganizationsRequest.create(payload)
})

export const create = (payload: IOrganization) => grpc.invoke({
    service,
    method : 'createOrganization',
    data   : CreateOrganizationRequest.create({ organization: payload })
})

export const update = (payload: IOrganization) => grpc.invoke({
    service,
    method : 'updateOrganization',
    data   : UpdateOrganizationRequest.create({ organization: payload })
})

export const destroy = (payload: IOrganization) => grpc.invoke({
    service,
    method : 'deleteOrganization',
    data   : DeleteOrganizationRequest.create(payload)
})

const actions = {
    get,
    list,
    create,
    update,
    destroy,
    select
}

const types = api.createActionTypes(resource)

export const {
    LIST_REQUESTED,    LIST_SUCCEEDED,    LIST_FAILED,
    GET_REQUESTED,     GET_SUCCEEDED,     GET_FAILED,
    CREATE_REQUESTED,  CREATE_SUCCEEDED,  CREATE_FAILED,
    UPDATE_REQUESTED,  UPDATE_SUCCEEDED,  UPDATE_FAILED,
    DESTROY_REQUESTED, DESTROY_SUCCEEDED, DESTROY_FAILED
} = types


//
//  REDUCER

const apiReducer = api.createReducer(resource, types)

const initialState = {
    ...api.initialState,
    current: null
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case SELECTED: {
            return {
                ...state,
                current: action.payload.name
            }
        }

        default:
            return apiReducer(state, action)
    }
}


//
//  SAGA

export function* saga() {
    yield fork(api.emitResourceActions, resource, types)
    yield takeEvery(app.INITIALIZED, decideOrganizationContext)
    yield takeLatest(SELECTED, setGRPCOrganizationMetadata)
    yield takeLatest(SELECTED, updateAddressWithOrganization)
}

function* decideOrganizationContext() {
    const organizations = yield _select(getAll)
    const currentlySelected = yield _select(getCurrent)

    if (isEmpty(organizations)) {
        yield put(list())
        const { success, failure } = yield race({
            success: take(LIST_SUCCEEDED),
            failure: take(LIST_FAILED)
        })
    }
    else {
        const params = yield _select(routing.getParams)
        const organizationFromAddress = yield _select(getByName(_get(params, 'org')))

        if (organizationFromAddress && !isEqual(organizationFromAddress, currentlySelected)) {
            yield put(select(organizationFromAddress))
        }
        else {
            if (currentlySelected) {
                yield put(select(currentlySelected))
            }
            else {
                const firstOrganizationAsDefault = head(values(organizations))
                if (firstOrganizationAsDefault) {
                    yield put(select(firstOrganizationAsDefault))
                }
            }
        }
    }
}

function setGRPCOrganizationMetadata(action: ActionType<typeof select>) {
    if (!action.payload.name) {
        return
    }

    grpc.setMetadata('organization', action.payload.name)
}

function* updateAddressWithOrganization(action: ActionType<typeof select>) {
    if (!action.payload.name) {
        return
    }

    const currentRoute = yield _select(routing.getCurrentRoute)
    const updatedURL = new URI(currentRoute.url)

    updatedURL.addSearch('org', action.payload.name)
    routing.replace(updatedURL.toString()) // eslint-disable-line lodash/prefer-lodash-method
}


//
//  SELECTORS

const selectors = api.createSelectors(resource)

export const { getState, getAll, countAll, getByName } = selectors

export const getCurrent = createSelector(
    [getState, getAll],
    (state, orgs) => find(orgs, { name: state.current }) || null
)
