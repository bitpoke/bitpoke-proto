import { ActionType, action as createAction } from 'typesafe-actions'
import { takeLatest, fork, put, take, race, select as _select, delay } from 'redux-saga/effects'
import { createSelector } from 'reselect'
import { find, head, replace, values as _values, get as _get, isEmpty, isEqual } from 'lodash'

import URI from 'urijs'

import { RootState, auth, api, grpc, routing, forms } from '../redux'

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

export type OrganizationName = string
export interface IOrganization extends presslabs.dashboard.organizations.v1.IOrganization {
    name: OrganizationName
}

export type Organization =
    presslabs.dashboard.organizations.v1.Organization

export type IOrganizationPayload =
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

export const { parseName, buildName } = api.createNameHelpers('orgs/:slug')


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

export const create = (payload: ICreateOrganizationRequest) => grpc.invoke({
    service,
    method : 'createOrganization',
    data   : CreateOrganizationRequest.create(payload)
})

export const update = (payload: IUpdateOrganizationRequest) => grpc.invoke({
    service,
    method : 'updateOrganization',
    data   : UpdateOrganizationRequest.create(payload)
})

export const destroy = (payload: IOrganizationPayload) => grpc.invoke({
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

const apiTypes = api.createActionTypes(resource)

export const {
    LIST_REQUESTED,    LIST_SUCCEEDED,    LIST_FAILED,
    GET_REQUESTED,     GET_SUCCEEDED,     GET_FAILED,
    CREATE_REQUESTED,  CREATE_SUCCEEDED,  CREATE_FAILED,
    UPDATE_REQUESTED,  UPDATE_SUCCEEDED,  UPDATE_FAILED,
    DESTROY_REQUESTED, DESTROY_SUCCEEDED, DESTROY_FAILED
} = apiTypes


//
//  REDUCER

const apiReducer = api.createReducer(resource, apiTypes)

const initialState: State = {
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
    yield fork(api.emitResourceActions, resource, apiTypes)
    yield takeLatest([
        routing.ROUTE_CHANGED,
        CREATE_SUCCEEDED,
        DESTROY_SUCCEEDED
    ], decideOrganizationContext)
    yield takeLatest(SELECTED, setGRPCOrganizationMetadata)
    yield fork(forms.takeEverySubmission, forms.Name.organization, handleFormSubmission)
}

const handleFormSubmission = api.createFormHandler(
    forms.Name.organization,
    api.Resource.organization,
    apiTypes,
    actions
)

function* decideOrganizationContext(): Iterable<any> {
    yield delay(50)

    const isAuthenticated = yield _select(auth.isAuthenticated)
    if (!isAuthenticated) {
        yield take(auth.ACCESS_GRANTED)
    }

    const organizations = yield _select(getAll)
    const currentlySelected = yield _select(getCurrent)

    if (isEmpty(organizations)) {
        yield put(list())
        const { success } = yield race({
            success: take(LIST_SUCCEEDED),
            failure: take(LIST_FAILED)
        })

        if (success && !api.isEmptyResponse(success)) {
            yield decideOrganizationContext()
            return
        }
    }
    else {
        const params = yield _select(routing.getParams)
        const name = buildName({ slug: _get(params, 'org') })
        const organizationFromAddress = name ? yield _select(getByName(name)) : null

        if (organizationFromAddress) {
            yield put(select(organizationFromAddress))

            // if (!isEqual(organizationFromAddress, currentlySelected)) {
            //     yield put(select(organizationFromAddress))
            // }
        }
        else {
            if (currentlySelected) {
                yield put(select(currentlySelected))
            }
            else {
                const firstOrganizationAsDefault = head(_values(organizations))
                if (firstOrganizationAsDefault) {
                    yield put(select(firstOrganizationAsDefault))
                }
            }
        }
    }
}

function* setGRPCOrganizationMetadata(action: ActionType<typeof select>) {
    yield put(grpc.setMetadata({
        key: 'organization',
        value: action.payload.name
    }))
}


//
//  SELECTORS

const selectors = api.createSelectors(resource)

export const getState:  (state: RootState) => State = selectors.getState
export const getAll:    (state: RootState) => api.ResourcesList<IOrganization> = selectors.getAll
export const countAll:  (state: RootState) => number = selectors.countAll
export const getByName: (name: OrganizationName) => (state: RootState) => IOrganization | null = selectors.getByName
export const getForURL: (url: routing.Path) => (state: RootState) => IOrganization | null = selectors.getForURL

export const getForCurrentURL = createSelector(
    [routing.getCurrentRoute, (state: RootState) => state],
    (currentRoute, state) => getForURL(currentRoute.url)(state)
)
export const getCurrent: (state: RootState) => IOrganization | null = createSelector(
    [getState, getAll],
    (state, orgs) => state.current
        ? find(orgs, { name: state.current }) || null
        : null
)
