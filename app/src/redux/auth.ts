import { ActionType, action as createAction } from 'typesafe-actions'
import { takeEvery, select, call, put } from 'redux-saga/effects'
import { channel as createChannel } from 'redux-saga'
import { User as Token, UserManager } from 'oidc-client'
import { createSelector } from 'reselect'
import axios from 'axios'

import { join, pick } from 'lodash'

import config from '../config'

import { RootState, app, routing, grpc } from '../redux'
import { watchChannel } from '../utils'

//
//  TYPES

export type State = Token | null
export type Actions = ActionType<typeof actions>

export type User = {
    id: string,
    name: string,
    email: string,
    isEmailVerified: boolean,
    nickname: string,
    avatarURL: string
} | null

//
//  ACTIONS

export const LOGIN_SUCCEEDED         = '@ auth / LOGIN_SUCCEEDED'
export const LOGIN_FAILED            = '@ auth / LOGIN_FAILED'
export const LOGOUT_REQUESTED        = '@ auth / LOGOUT_REQUESTED'
export const TOKEN_REFRESH_REQUESTED = '@ auth / TOKEN_REFRESH_REQUESTED'
export const ACCESS_GRANTED          = '@ auth / ACCESS_GRANTED'

export const loginSuccess = (token: Token) => createAction(LOGIN_SUCCEEDED, token)
export const loginFailure = (error: any) => createAction(LOGIN_FAILED, error)
export const logout = () => createAction(LOGOUT_REQUESTED)
export const refreshToken = () => createAction(TOKEN_REFRESH_REQUESTED)
export const grantAccess = () => createAction(ACCESS_GRANTED)

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

function* ensureAuthentication(action: ActionType<typeof routing.updateRoute>) {
    const userIsAuthenticated = yield select(isAuthenticated)
    const route = yield select(routing.getCurrentRoute)
    if (userIsAuthenticated) {
        yield call(setGRPCAuthorizationMetadata)
        // yield put(grantAccess())
        provider.startSilentRenew()
        return
    }

    if (route && hasAuthenticationPayload(route.path)) {
        provider.signinRedirectCallback()
        return
    }

    if (!userIsAuthenticated) {
        provider.signinRedirect()
    }

    return
}

function handleAuthenticationIfRequired(action: ActionType<typeof routing.updateRoute>) {
    const { pathname } = action.payload
    if (hasAuthenticationPayload(pathname)) {
        provider.signinRedirectCallback()
    }
}

function handleTokenRefresh(action: ActionType<typeof refreshToken>) {
    provider.signinSilent()
}

function redirectToDashboard() {
    routing.push(routing.routeFor('dashboard'))
}

function* setGRPCAuthorizationMetadata() {
    const token = yield select(getToken)
    grpc.setMetadata('Authorization', token)
}


//
//   HELPERS and UTILITIES

const provider = new UserManager({
    authority            : config.REACT_APP_OIDC_ISSUER,
    client_id            : config.REACT_APP_OIDC_CLIENT_ID,
    redirect_uri         : window.location.href,
    silent_redirect_uri  : window.location.href,
    response_type        : 'token id_token',
    scope                : 'openid email profile',
    automaticSilentRenew : true
})

provider.events.addUserLoaded((token) => {
    token && !token.expired
        ? channel.put(loginSuccess(token))
        : channel.put(loginFailure(token))
})
provider.events.addUserUnloaded(() => channel.put(logout()))
provider.events.addUserSignedOut(() => channel.put(logout()))
provider.events.addAccessTokenExpired(() => channel.put(logout()))

function hasAuthenticationPayload(path: string) {
    return /access_token|id_token|error/.test(path)
}

function tokenIsValid(token: State) {
    if (!token || !token.expires_at) {
        return false
    }
    return (Date.now() / 1000) < token.expires_at
}

function normalizeProfile(token: Token) {
    if (!tokenIsValid(token)) {
        return null
    }

    return {
        id: token.profile.sub,
        isEmailVerified: token.profile.email_verified,
        avatarURL: token.profile.picture,
        ...pick(token.profile, ['email', 'name', 'nickname'])
    }
}

//
//  SELECTORS

export const getState = (state: RootState): State => state.auth
export const getToken = createSelector(
    getState,
    (state) => state ? join([state.token_type, state.id_token], ' ') : null
)
export const getCurrentUser = createSelector(
    getState,
    (state) => state ? normalizeProfile(state) : null
)
export const isAuthenticated = createSelector(
    getState,
    (state) => tokenIsValid(state)
)
