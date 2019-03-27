import { ActionType } from 'typesafe-actions'
import { createSelector } from 'reselect'

import { RootState } from '../redux'


//
//  TYPES

export type State = {
    isInitialized: boolean
}
export type Actions = ActionType<typeof actions>


//
//  ACTIONS

export const INITIALIZED = '@ app / INITIALIZED'
export const initialize = () => ({ type: INITIALIZED })

const actions = {
    initialize
}


//
//  REDUCER

const initialState: State = {
    isInitialized: false
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case INITIALIZED:
            return {
                ...state,
                isInitialized: true
            }

        default:
            return state
    }
}

//
//  SELECTORS

export const getState = (state: RootState): State => state.app
export const isInitialized = createSelector(
    getState,
    (state) => state.isInitialized
)
