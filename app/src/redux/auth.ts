import { ActionType, action as createAction } from 'typesafe-actions'
import { takeEvery, select } from 'redux-saga/effects'
import { channel as createChannel } from 'redux-saga'
import { Auth0DecodedHash, Auth0Error, WebAuth } from 'auth0-js'
import { createSelector } from 'reselect'

import config from '../config'

import { join, pick } from 'lodash'

import { RootState, app, routing } from '../redux'
import { watchChannel } from '../utils'

//
//  TYPES

export type State = Auth0DecodedHash | null
export type Actions = ActionType<typeof actions>

type TokenPayload = {
    exp: number
}

export type User = {
    id: string,
    name: string,
    email: string,
    isEmailVerified: boolean,
    nickname: string,
    avatarURL: string
}

//
//  ACTIONS

export const LOGIN_SUCCEEDED         = '@ auth / LOGIN_SUCCEEDED'
export const LOGIN_FAILED            = '@ auth / LOGIN_FAILED'
export const LOGOUT_REQUESTED        = '@ auth / LOGOUT_REQUESTED'
export const TOKEN_REFRESH_REQUESTED = '@ auth / TOKEN_REFRESH_REQUESTED'

export const loginSuccess = (hash: Auth0DecodedHash) => createAction(LOGIN_SUCCEEDED, hash)
export const loginFailure = (error: Auth0Error) => createAction(LOGIN_FAILED, error)
export const logout = () => createAction(LOGOUT_REQUESTED)
export const refreshToken = () => createAction(TOKEN_REFRESH_REQUESTED)

const actions = {
    logout,
    loginSuccess,
    loginFailure,
    refreshToken
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
    yield takeEvery(routing.ROUTE_CHANGED, ensureAuthentication)
    yield takeEvery([LOGIN_SUCCEEDED, LOGIN_FAILED, LOGOUT_REQUESTED], redirectToDashboard)
    yield takeEvery(TOKEN_REFRESH_REQUESTED, handleTokenRefresh)
    yield watchChannel(channel)
}

function* ensureAuthentication(action: ActionType<typeof app.initialize>) {
    const userIsAuthenticated = yield select(isAuthenticated)
    const route = yield select(routing.getCurrentRoute)

    if (userIsAuthenticated) {
        return
    }

    if (route && hasAuthenticationPayload(route.path)) {
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

function handleTokenRefresh(action: ActionType<typeof refreshToken>) {
    provider.checkSession({}, handleTokenResponse)
}

function handleTokenResponse(err: Auth0Error | null, authResult: Auth0DecodedHash) {
    if (authResult && authResult.accessToken && authResult.idToken) {
        channel.put(loginSuccess(authResult))
    } else if (err) {
        channel.put(loginFailure(err))
    }
}

function redirectToDashboard() {
    routing.push(routing.routeFor('dashboard'))
}


//
//   HELPERS and UTILITIES

const provider = new WebAuth({
    domain       : config.REACT_APP_AUTH0_DOMAIN || '{DOMAIN}',
    clientID     : config.REACT_APP_AUTH0_CLIENT_ID || '{CLIENT_ID}',
    redirectUri  : config.REACT_APP_AUTH0_CALLBACK_URL || 'http://localhost:3000/',
    audience     : `https://${config.REACT_APP_AUTH0_DOMAIN}/userinfo`,
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
export const getCurrentUser = createSelector(
    getTokenPayload,
    (token): User => ({
        id: token.sub,
        isEmailVerified: token.email_verified,
        avatarURL: token.picture,
        ...pick(token, ['email', 'name', 'nickname'])
    })
)
export const isAuthenticated = createSelector(
    getTokenPayload,
    (token) => tokenIsValid(token)
)
