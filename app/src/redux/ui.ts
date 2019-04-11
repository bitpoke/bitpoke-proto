import { ActionType, action as createAction } from 'typesafe-actions'
import { channel as createChannel } from 'redux-saga'
import { takeEvery } from 'redux-saga/effects'
import { createSelector } from 'reselect'
import { Position, Toaster, Intent, IconName } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState } from '../redux'
import { watchChannel } from '../utils'


const DEFAULT_TIMEOUT = 3000

//
//  TYPES

export type State = {}
export type Actions = ActionType<typeof actions>

export type Toast = {
    message: string,
    icon?: IconName,
    intent?: Intent,
    isPersistent?: boolean
}


//
//  ACTIONS

export const TOAST_DISPLAYED = '@ ui / TOAST_DISPLAYED'
export const TOAST_CLOSED    = '@ ui / TOAST_CLOSED'

export const showToast = (payload: Toast) => createAction(TOAST_DISPLAYED, payload)
export const closeToast = () => createAction(TOAST_CLOSED)

const actions = {
    showToast
}


//
//  REDUCER

const initialState: State = {}

export function reducer(state: State = initialState, action: Actions) {
    return state
}


//
//  SAGA

const channel = createChannel()
const toaster = Toaster.create({ position: Position.TOP })

export function* saga() {
    yield takeEvery(TOAST_DISPLAYED, dispatchToToaster)
    yield watchChannel(channel)
}

function* dispatchToToaster(action: ActionType<typeof showToast>): Iterator<any> {
    const toast = action.payload
    const timeout = toast.isPersistent ? 0 : get(toast, 'timeout', DEFAULT_TIMEOUT)

    toaster.show({
        timeout,
        ...toast,
        onDismiss: () => channel.put(closeToast())
    })
}


//
//  SELECTORS

export const getState = (state: RootState): State => state.ui
