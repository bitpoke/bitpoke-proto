import { ActionType, action as createAction } from 'typesafe-actions'
import { takeEvery, select, put, call } from 'redux-saga/effects'
import { channel as createChannel } from 'redux-saga'
import { User as Token, UserManager } from 'oidc-client'
import { createSelector } from 'reselect'

import { join, pick } from 'lodash'

import config from '../config'

import { RootState, app, routing, grpc } from '../redux'
import { watchChannel } from '../utils'

//
//  TYPES

export type State = {
    token: Token | null
    isAccessGranted: boolean
}

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
export const logout       = () => createAction(LOGOUT_REQUESTED)
export const refreshToken = () => createAction(TOKEN_REFRESH_REQUESTED)
export const grantAccess  = () => createAction(ACCESS_GRANTED)

const actions = {
    logout,
    loginSuccess,
    loginFailure,
    refreshToken,
    grantAccess
}


//
//  REDUCER

const initialState: State = {
    token: null,
    isAccessGranted: false
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case LOGIN_SUCCEEDED: {
            return {
                ...state,
                token: action.payload
            }
        }

        case LOGIN_FAILED:
        case LOGOUT_REQUESTED: {
            return null
        }

        case ACCESS_GRANTED: {
            return {
                ...state,
                isAccessGranted: true
            }
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
    yield takeEvery([
        app.INITIALIZED,
        LOGIN_SUCCEEDED
    ], grantAccessIfRequired)
    yield watchChannel(channel)
}

function* ensureAuthentication(action: ActionType<typeof routing.updateRoute>) {
    const hasToken = yield select(hasValidToken)
    const route = yield select(routing.getCurrentRoute)
    if (hasToken) {
        provider.startSilentRenew()
        return
    }

    if (route && hasAuthenticationPayload(route.path)) {
        provider.signinRedirectCallback()
        return
    }

    if (!hasToken) {
        provider.signinRedirect()
    }

    return
}

function* grantAccessIfRequired() {
    const hasToken = yield select(hasValidToken)
    if (hasToken) {
        yield call(setGRPCAuthorizationMetadata)
        yield put(grantAccess())
    }
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

function* redirectToDashboard() {
    yield put(routing.push(routing.routeFor('dashboard')))
}

function* setGRPCAuthorizationMetadata() {
    const tokenPayload = yield select(getTokenPayload)
    yield put(grpc.setMetadata({
        key: 'Authorization',
        value: tokenPayload
    }))
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

function tokenIsValid(token: Token | null) {
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
    (state) => state.token
)
export const getTokenPayload = createSelector(
    getToken,
    (token) => token ? join([token.token_type, token.id_token], ' ') : null
)
export const getCurrentUser = createSelector(
    getToken,
    (token) => token ? normalizeProfile(token) : null
)
export const hasValidToken = createSelector(
    getToken,
    (token) => token ? tokenIsValid(token) : false
)
export const isAuthenticated = createSelector(
    [hasValidToken, getState],
    (hasToken, state) => hasToken && state.isAccessGranted
)
