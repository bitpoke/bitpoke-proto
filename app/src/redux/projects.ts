import { ActionType, createAsyncAction, action as createAction, isOfType } from 'typesafe-actions'
import { takeEvery, takeLatest, fork, put, take, call, race, select } from 'redux-saga/effects'
import { SagaIterator, channel as createChannel } from 'redux-saga'
import { createSelector } from 'reselect'

import { reduce, pickBy, get as _get, join, noop, snakeCase, toUpper, includes, isEmpty } from 'lodash'

import { RootState, api, grpc, forms, organizations, routing, toasts, wizards } from '../redux'

import { presslabs } from '@presslabs/dashboard-proto'

import { Intent } from '@blueprintjs/core'

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

export type ProjectName = string
export interface IProject extends presslabs.dashboard.projects.v1.IProject {
    name: ProjectName
}

export type Project =
    presslabs.dashboard.projects.v1.Project

export type IProjectPayload =
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

export const create = (payload: IProjectPayload) => grpc.invoke({
    service,
    method : 'createProject',
    data   : CreateProjectRequest.create({ project: payload })
})

export const update = (payload: IProjectPayload) => grpc.invoke({
    service,
    method : 'updateProject',
    data   : UpdateProjectRequest.create({ project: payload })
})

export const destroy = (payload: IProjectPayload) => grpc.invoke({
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

const apiTypes = api.createActionTypes(api.Resource.project)

export const {
    LIST_REQUESTED,    LIST_SUCCEEDED,    LIST_FAILED,
    GET_REQUESTED,     GET_SUCCEEDED,     GET_FAILED,
    CREATE_REQUESTED,  CREATE_SUCCEEDED,  CREATE_FAILED,
    UPDATE_REQUESTED,  UPDATE_SUCCEEDED,  UPDATE_FAILED,
    DESTROY_REQUESTED, DESTROY_SUCCEEDED, DESTROY_FAILED
} = apiTypes


//
//  REDUCER

const apiReducer = api.createReducer(api.Resource.project, apiTypes)

export function reducer(state: State = api.initialState, action: Actions) {
    return apiReducer(state, action)
}


//
//  SAGA

export function* saga() {
    yield fork(api.emitResourceActions, api.Resource.project, apiTypes)
    yield takeLatest(organizations.SELECTED, fetchAll)
    yield forms.takeEverySubmission(forms.Name.project, handleFormSubmission)
}

function* fetchAll() {
    yield put(list())
    const projects = yield select(getAll)

    if (isEmpty(projects)) {
        const { success, failure } = yield race({
            success: take(LIST_SUCCEEDED),
            failure: take(LIST_FAILED)
        })

        if (success && api.isEmptyResponse(success)) {
            // yield put(wizards.startFlow('onboarding'))
        }
    }
}

function* handleFormSubmission(action: forms.Actions) {
    const { resolve, reject, values } = action.payload

    if (api.isNewEntry(values)) {
        yield put(create(values))

        const { success, failure } = yield race({
            success : take(CREATE_SUCCEEDED),
            failure : take(CREATE_FAILED)
        })

        if (success) {
            yield call(resolve)
            yield put(forms.reset(forms.Name.project))
            yield put(routing.push(routing.routeFor('dashboard')))
            toasts.show({
                intent: Intent.SUCCESS,
                message: 'Project created'
            })
        }
        else {
            yield call(reject, new forms.SubmissionError())
            toasts.show({
                intent: Intent.DANGER,
                message: 'Failed to create project'
            })
        }
    }
    else {
        yield put(update(values))

        const { success, failure } = yield race({
            success : take(UPDATE_SUCCEEDED),
            failure : take(UPDATE_FAILED)
        })

        if (success) {
            yield call(resolve)
            yield put(forms.reset(forms.Name.project))
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

const selectors = api.createSelectors(api.Resource.project)

export const getState: (state: RootState) => State = selectors.getState
export const getAll: (state: RootState) => api.ResourcesList<IProject> = selectors.getAll
export const countAll: (state: RootState) => number = selectors.countAll
export const getByName: (name: ProjectName) => (state: RootState) => IProject | null = selectors.getByName

export const getForCurrentOrganization = createSelector(
    [organizations.getCurrent, getAll],
    (currentOranization, projects) => currentOranization
        ? pickBy(projects, { organization: currentOranization.name })
        : {}
)
