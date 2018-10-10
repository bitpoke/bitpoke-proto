import { ActionType, action as createAction } from 'typesafe-actions'
import { takeEvery, select } from 'redux-saga/effects'
import { channel as createChannel } from 'redux-saga'
import { grpc } from 'grpc-web-client'
import { createSelector } from 'reselect'

import config from '../config'
import { watchChannel } from '../utils'
import { RootState, auth } from '../redux'

import { ListRequest, Project } from '../proto/presslabs/dashboard/core/v1/project_pb'
import { Projects } from '../proto/presslabs/dashboard/core/v1/project_pb_service'

const host: string = config.REACT_API_URL || 'http://localhost:8080'


//
//  TYPES

export type State = {
    readonly entries: ProjectsList
}
export type Actions = ActionType<typeof actions>
export { Project }

export type ProjectsList = {
    [id: string]: Project;
}

//
//  ACTIONS

export const LIST_REQUESTED = '@ projects / LIST_REQUESTED'
export const LIST_SUCCEEDED = '@ projects / LIST_SUCCEEDED'

export const RECEIVED       = '@ projects / RECEIVED'

export const list = () => createAction(LIST_REQUESTED)
export const receive = (entry: Project) => createAction(RECEIVED, entry)

const listRequest = {
    service: Projects.List,
    request: new ListRequest()
}

const actions = {
    list,
    receive
}


//
//  REDUCER

const initialState: State = {
    entries: {}
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case RECEIVED: {
            const project = action.payload
            return {
                ...state,
                entries: {
                    ...state.entries,
                    [project.getId()]: project
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
    yield takeEvery(LIST_REQUESTED, performRequest)
    yield watchChannel(channel)
}

function* performRequest(action: ActionType<typeof list>) {
    const authorization = yield select(auth.getAuthorizationHeader)
    grpc.invoke(listRequest.service, {
        host,
        request: listRequest.request,
        metadata: { authorization },
        onMessage: (response: Project) =>
            channel.put(receive(response)),
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
    (state: State): ProjectsList => state.entries
)
