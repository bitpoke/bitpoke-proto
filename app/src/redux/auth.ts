import { ActionType, action as createAction } from 'typesafe-actions'
import { takeEvery, fork, put, select, take, call } from 'redux-saga/effects'
import { Channel, SagaIterator, channel as createChannel } from 'redux-saga'
import auth0 from 'auth0-js'
import { createSelector } from 'reselect'

import { get, join } from 'lodash'

import { AnyAction, RootState, app, routing } from '../redux'

//
//  TYPES

export type State = auth0.Auth0DecodedHash | null
export type Actions = ActionType<typeof actions>

type TokenPayload = {
    exp: number
}

//
//  ACTIONS

export const LOGIN_SUCCEEDED  = '@ auth / LOGIN_SUCCEEDED'
export const LOGIN_FAILED     = '@ auth / LOGIN_FAILED'
export const LOGOUT_REQUESTED = '@ auth / LOGOUT_REQUESTED'

export const loginSuccess = (hash: auth0.Auth0DecodedHash) => createAction(LOGIN_SUCCEEDED, hash)
export const loginFailure = (error: auth0.Auth0Error) => createAction(LOGIN_FAILED, error)
export const logout = () => createAction(LOGOUT_REQUESTED)

const actions = {
    logout,
    loginSuccess,
    loginFailure
}


//
//  REDUCER

const initialState: State = null
export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case LOGIN_SUCCEEDED: {
            return action.payload
        }

        case LOGIN_FAILED:
        case LOGOUT_REQUESTED: {
            return null
        }

        default:
            return state
    }
    return state
}


//
//  SAGA

const channel = createChannel()

export function* saga() {
    yield takeEvery(routing.ROUTE_CHANGED, handleAuthenticationIfRequired)
    yield takeEvery(app.INITIALIZED, ensureAuthentication)
    yield takeEvery([LOGIN_SUCCEEDED, LOGIN_FAILED], redirectToDashboard)
    yield watchChannel(channel)
}

function* ensureAuthentication(action: ActionType<typeof app.initialize>) {
    const userIsAuthenticated = yield select(isAuthenticated)
    const route = yield select(routing.getCurrentRoute)

    if (hasAuthenticationPayload(route.path)) {
        return
    }

    if (!userIsAuthenticated) {
        provider.authorize()
    }

    return
}

function handleAuthenticationIfRequired(action: ActionType<typeof routing.updateRoute>) {
    if (hasAuthenticationPayload(action.payload)) {
        provider.parseHash(handleTokenResponse)
    }
}

function handleTokenResponse(err: auth0.Auth0Error | null, authResult: auth0.Auth0DecodedHash) {
    if (authResult && authResult.accessToken && authResult.idToken) {
        channel.put(loginSuccess(authResult))
    } else if (err) {
        channel.put(loginFailure(err))
    }
}

function redirectToDashboard() {
    routing.push(routing.routeFor('dashboard'))
}

function* watchChannel(actionChannel: Channel<{}>) {
    while (true) {
        yield put(yield take(actionChannel))
    }
}

//
//   HELPERS and UTILITIES

const provider = new auth0.WebAuth({
    domain       : process.env.REACT_APP_AUTH0_DOMAIN || '{DOMAIN}',
    clientID     : process.env.REACT_APP_AUTH0_CLIENT_ID || '{CLIENT_ID}',
    redirectUri  : process.env.REACT_APP_AUTH0_CALLBACK_URL || 'http://localhost:3000/',
    audience     : `https://${process.env.REACT_APP_AUTH0_DOMAIN}/userinfo`,
    responseType : 'token id_token',
    scope        : 'openid email profile'
})

function hasAuthenticationPayload(path: string) {
    return /access_token|id_token|error/.test(path)
}

function tokenIsValid(token: TokenPayload) {
    return (token && token.exp && (Date.now() / 1000) < token.exp) || false
}


//
//  SELECTORS

export const getState = (state: RootState): State => state.auth
export const getTokenPayload = createSelector(
    getState,
    (state) => state ? state.idTokenPayload : null
)
export const getAuthorizationHeader = createSelector(
    getState,
    (state) => state ? join([state.tokenType, state.idToken], ' ') : null
)
export const isAuthenticated = createSelector(
    getTokenPayload,
    (token) => tokenIsValid(token)
)
