import { ActionType, createAsyncAction, action as createAction } from 'typesafe-actions'
import { takeEvery, put, select, take, call } from 'redux-saga/effects'
import { grpc } from 'grpc-web-client'
import { createSelector } from 'reselect'

import { RootState, auth } from '../redux'

import { ListProjectsRequest, Project } from '../proto/projects/v1/project_pb'
import { Projects } from '../proto/projects/v1/project_pb_service'

const host: string = process.env.REACT_API_URL || 'http://localhost:9090'


//
//  TYPES

export type State = {
    readonly entries: Project[]
}
export type Actions = ActionType<typeof actions>


//
//  ACTIONS

export const LIST_REQUESTED = '@ projects / LIST_REQUESTED'
export const LIST_SUCCEEDED = '@ projects / LIST_SUCCEEDED'
export const LIST_FAILED    = '@ projects / LIST_FAILED'

export const list = () => createAction(LIST_REQUESTED)

const listRequest = {
    service: Projects.ListProjects,
    request: new ListProjectsRequest()
}

const actions = {
    list
}


//
//  REDUCER

const initialState: State = {
    entries: []
}

export function reducer(state: State = initialState, action: Actions) {
    return state
}


//
//  SAGA

export function* saga() {
    yield takeEvery(LIST_REQUESTED, performRequest)
}

function* performRequest(action: ActionType<typeof list>) {
    const authorization = yield select(auth.getAuthorizationHeader)
    grpc.invoke(listRequest.service, {
        host,
        request: listRequest.request,
        metadata: { authorization },
        onMessage: (response: Project) =>
            console.log(response.toObject()),
        onEnd: () =>
            console.log('ONEND!'),
        debug: true
    })
}


//
//  SELECTORS

export const getState = (state: RootState): State => state.projects
export const getAll = createSelector(
    getState,
    (state: State) => []
)
