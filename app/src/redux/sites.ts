import { ActionType } from 'typesafe-actions'
import { fork, put, take, call, race } from 'redux-saga/effects'
import { createSelector } from 'reselect'

import { pickBy, get as _get, startsWith } from 'lodash'

import { RootState, api, grpc, forms, projects, routing, toasts } from '../redux'

import { presslabs } from '@presslabs/dashboard-proto'

import { Intent } from '@blueprintjs/core'

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

export type SiteName = string
export interface ISite extends presslabs.dashboard.sites.v1.ISite {
    name: SiteName
}

export type Site =
    presslabs.dashboard.sites.v1.Site

export type ISitePayload =
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

export const create = (payload: ICreateSiteRequest) => grpc.invoke({
    service,
    method : 'createSite',
    data   : CreateSiteRequest.create(payload)
})

export const update = (payload: IUpdateSiteRequest) => grpc.invoke({
    service,
    method : 'updateSite',
    data   : UpdateSiteRequest.create(payload)
})

export const destroy = (payload: ISitePayload) => grpc.invoke({
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

const apiTypes = api.createActionTypes(api.Resource.site)

export const {
    LIST_REQUESTED,    LIST_SUCCEEDED,    LIST_FAILED,
    GET_REQUESTED,     GET_SUCCEEDED,     GET_FAILED,
    CREATE_REQUESTED,  CREATE_SUCCEEDED,  CREATE_FAILED,
    UPDATE_REQUESTED,  UPDATE_SUCCEEDED,  UPDATE_FAILED,
    DESTROY_REQUESTED, DESTROY_SUCCEEDED, DESTROY_FAILED
} = apiTypes


//
//  REDUCER

const apiReducer = api.createReducer(api.Resource.site, apiTypes)

export function reducer(state: State = api.initialState, action: Actions) {
    return apiReducer(state, action)
}


//
//  SAGA

export function* saga() {
    yield fork(api.emitResourceActions, api.Resource.site, apiTypes)
    yield forms.takeEverySubmission(forms.Name.site, handleFormSubmission)
}

function* handleFormSubmission(action: forms.Actions) {
    const { resolve, reject, values } = action.payload
    const entry = _get(values, api.Resource.site)

    if (api.isNewEntry(entry)) {
        yield put(create(values))

        const { success } = yield race({
            success : take(CREATE_SUCCEEDED),
            failure : take(CREATE_FAILED)
        })

        if (success) {
            yield call(resolve)
            yield put(forms.reset(forms.Name.site))
            yield put(routing.push(routing.routeFor('dashboard')))
            toasts.show({
                intent: Intent.SUCCESS,
                message: 'Site created'
            })
        }
        else {
            yield call(reject, new forms.SubmissionError())
            toasts.show({
                intent: Intent.DANGER,
                message: 'Failed to create site'
            })
        }
    }
    else {
        yield put(update(values))

        const { success } = yield race({
            success : take(UPDATE_SUCCEEDED),
            failure : take(UPDATE_FAILED)
        })

        if (success) {
            yield call(resolve)
            yield put(forms.reset(forms.Name.site))
            yield put(routing.push(routing.routeFor('dashboard')))
            toasts.show({
                intent: Intent.SUCCESS,
                message: 'Project updated'
            })
        }
        else {
            yield call(reject, new forms.SubmissionError())
            toasts.show({
                intent: Intent.DANGER,
                message: 'Failed to update the project'
            })
        }
    }
}

//
//  SELECTORS

const selectors = api.createSelectors(api.Resource.site)

export const getState:  (state: RootState) => State = selectors.getState
export const getAll:    (state: RootState) => api.ResourcesList<ISite> = selectors.getAll
export const countAll:  (state: RootState) => number = selectors.countAll
export const getByName: (name: SiteName) => (state: RootState) => ISite | null = selectors.getByName
export const getForURL: (url: routing.Path) => (state: RootState) => ISite | null = selectors.getForURL

export const getForProject = (project: projects.ProjectName) => createSelector(
    getAll,
    (sites) => pickBy(sites, ({ name }) => startsWith(name, project))
)
