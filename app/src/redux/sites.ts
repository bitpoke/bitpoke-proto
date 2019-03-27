import { ActionType, createAsyncAction, action as createAction, isOfType } from 'typesafe-actions'
import { takeEvery, takeLatest, fork, put, take, call } from 'redux-saga/effects'
import { SagaIterator, channel as createChannel } from 'redux-saga'
import { createSelector } from 'reselect'

import { reduce, pickBy, get as _get, join, noop, snakeCase, toUpper, includes } from 'lodash'

import { RootState, api, grpc } from '../redux'

import { presslabs } from '@presslabs/dashboard-proto'

const {
    Site,
    SitesService,
    ListSitesRequest,
    ListSitesResponse,
    GetSiteRequest,
    CreateSiteRequest,
    UpdateSiteRequest,
    DeleteSiteRequest
} = presslabs.dashboard.sites.v1

export {
    Site,
    SitesService,
    ListSitesRequest,
    ListSitesResponse,
    GetSiteRequest,
    CreateSiteRequest,
    UpdateSiteRequest,
    DeleteSiteRequest
}


//
//  TYPES

export type Site =
    presslabs.dashboard.sites.v1.Site

export type ISite =
    presslabs.dashboard.sites.v1.ISite

export type ListSitesResponse =
    presslabs.dashboard.sites.v1.ListSitesRequest

export type IListSitesRequest =
    presslabs.dashboard.sites.v1.IListSitesRequest

export type IGetSiteRequest =
    presslabs.dashboard.sites.v1.IGetSiteRequest

export type ICreateSiteRequest =
    presslabs.dashboard.sites.v1.ICreateSiteRequest

export type IUpdateSiteRequest =
    presslabs.dashboard.sites.v1.IUpdateSiteRequest

export type IDeleteSiteRequest =
    presslabs.dashboard.sites.v1.IDeleteSiteRequest

export type IListSitesResponse =
    presslabs.dashboard.sites.v1.IListSitesResponse

export type State = api.State<ISite>
export type Actions = ActionType<typeof actions>

const resource = api.Resource.site

const service = SitesService.create(
    grpc.createTransport('presslabs.dashboard.sites.v1.SitesService')
)


//
//  ACTIONS

export const get = (payload: IGetSiteRequest) => grpc.invoke({
    service,
    method : 'getSite',
    data   : GetSiteRequest.create(payload)
})

export const list = (payload?: IListSitesRequest) => grpc.invoke({
    service,
    method : 'listSites',
    data   : ListSitesRequest.create(payload)
})

export const create = (payload: ISite) => grpc.invoke({
    service,
    method : 'createSite',
    data   : CreateSiteRequest.create({ site: payload })
})

export const update = (payload: ISite) => grpc.invoke({
    service,
    method : 'updateSite',
    data   : UpdateSiteRequest.create({ site: payload })
})

export const destroy = (payload: ISite) => grpc.invoke({
    service,
    method : 'deleteSite',
    data   : DeleteSiteRequest.create(payload)
})

const actions = {
    get,
    list,
    create,
    update,
    destroy
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

export function reducer(state: State = api.initialState, action: Actions) {
    return apiReducer(state, action)
}


//
//  SAGA

export function* saga() {
    yield fork(api.emitResourceActions, resource, types)
}


//
//  SELECTORS

const selectors = api.createSelectors(resource)

export const { getState, getAll, countAll, getByName } = selectors

export const getForProject = (project) => createSelector(
    getAll,
    (sites) => pickBy(sites, { project })
)
