import { ActionType, createAsyncAction, action as createAction, isOfType } from 'typesafe-actions'
import { takeEvery, takeLatest, fork, put, take, call } from 'redux-saga/effects'
import { SagaIterator, channel as createChannel } from 'redux-saga'
import { createSelector } from 'reselect'

import { reduce, pickBy, get as _get, join, noop, snakeCase, toUpper, includes } from 'lodash'

import { RootState, api, grpc, organizations } from '../redux'

import { presslabs } from '@presslabs/dashboard-proto'

const {
    Project,
    ProjectsService,
    ListProjectsRequest,
    ListProjectsResponse,
    GetProjectRequest,
    CreateProjectRequest,
    UpdateProjectRequest,
    DeleteProjectRequest
} = presslabs.dashboard.projects.v1

export {
    Project,
    ProjectsService,
    ListProjectsRequest,
    ListProjectsResponse,
    GetProjectRequest,
    CreateProjectRequest,
    UpdateProjectRequest,
    DeleteProjectRequest
}


//
//  TYPES

export type Project =
    presslabs.dashboard.projects.v1.Project

export type IProject =
    presslabs.dashboard.projects.v1.IProject

export type ListProjectsResponse =
    presslabs.dashboard.projects.v1.ListProjectsRequest

export type IListProjectsRequest =
    presslabs.dashboard.projects.v1.IListProjectsRequest

export type IGetProjectRequest =
    presslabs.dashboard.projects.v1.IGetProjectRequest

export type ICreateProjectRequest =
    presslabs.dashboard.projects.v1.ICreateProjectRequest

export type IUpdateProjectRequest =
    presslabs.dashboard.projects.v1.IUpdateProjectRequest

export type IDeleteProjectRequest =
    presslabs.dashboard.projects.v1.IDeleteProjectRequest

export type IListProjectsResponse =
    presslabs.dashboard.projects.v1.IListProjectsResponse

export type State = api.State<IProject>
export type Actions = ActionType<typeof actions>

const resource = api.Resource.project

const service = ProjectsService.create(
    grpc.createTransport('presslabs.dashboard.projects.v1.ProjectsService')
)


//
//  ACTIONS

export const get = (payload: IGetProjectRequest) => grpc.invoke({
    service,
    method : 'getProject',
    data   : GetProjectRequest.create(payload)
})

export const list = (payload?: IListProjectsRequest) => grpc.invoke({
    service,
    method : 'listProjects',
    data   : ListProjectsRequest.create(payload)
})

export const create = (payload: IProject) => grpc.invoke({
    service,
    method : 'createProject',
    data   : CreateProjectRequest.create({ project: payload })
})

export const update = (payload: IProject) => grpc.invoke({
    service,
    method : 'updateProject',
    data   : UpdateProjectRequest.create({ project: payload })
})

export const destroy = (payload: IProject) => grpc.invoke({
    service,
    method : 'deleteProject',
    data   : DeleteProjectRequest.create(payload)
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
    yield takeLatest(organizations.SELECTED, fetchProjects)
}

function* fetchProjects() {
    yield put(list())
}


//
//  SELECTORS

const selectors = api.createSelectors(resource)

export const { getState, getAll, countAll, getByName } = selectors

export const getForCurrentOrganization = createSelector(
    [organizations.getCurrent, getAll],
    (currentOranization, projects) => currentOranization
        ? pickBy(projects, { organization: _get(currentOranization, 'name', 'test') })
        : {}
)
