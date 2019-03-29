import * as React from 'react'
import { connect } from 'react-redux'
import { createHashHistory, createMemoryHistory, Location, Path } from 'history'
import pathToRegexp, { Key, compile } from 'path-to-regexp'
import { matchPath } from 'react-router'
import { takeEvery, put } from 'redux-saga/effects'
import { channel as createChannel } from 'redux-saga'
import { ActionType, action as createAction, isOfType } from 'typesafe-actions'
import { createSelector } from 'reselect'

import URI from 'urijs'

import { map, filter, findKey, omit, omitBy, get, head, has, isEmpty, isString } from 'lodash'

import { RootState, app } from '../redux'
import { watchChannel } from '../utils'


//
//  TYPES

export type State = Route
export type Actions = ActionType<typeof actions>
export type Params = object

export type Route = {
    path   : string,
    url    : string,
    params : Params,
    key?   : string
}

//
//  ROUTES

export const ROUTE_MAP = {
    dashboard: {
        path      : '/',
        component : 'DashboardContainer'
    },
    onboarding: {
        path      : '/onboarding/:step?',
        component : 'OnboardingContainer'
    }
}

const INITIAL_ROUTE = {
    path   : '/',
    url    : '/',
    key    : findKey(ROUTE_MAP, { path: '/' }),
    params : {}
} as Route



//
//  ACTIONS

export const PUSH_REQUESTED    = '@ routing / PUSH_REQUESTED'
export const REPLACE_REQUESTED = '@ routing / REPLACE_REQUESTED'
export const ROUTE_CHANGED     = '@ routing / ROUTE_CHANGED'

export const push = (path: Path) => createAction(PUSH_REQUESTED, path)
export const replace = (path: Path) => createAction(REPLACE_REQUESTED, path)
export const updateRoute = (location: Location<any>) => createAction(ROUTE_CHANGED, location)

const actions = {
    push, replace, updateRoute
}

//
//  REDUCER

const initialState: State = INITIAL_ROUTE

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case ROUTE_CHANGED: {
            return matchRoute(action.payload)
        }

        default:
            return state
    }
}


//
//   SAGA

const channel = createChannel()

export function* saga() {
    yield takeEvery(app.INITIALIZED, bootstrap)
    yield takeEvery([PUSH_REQUESTED, REPLACE_REQUESTED], dispatchToHistory)
    yield watchChannel(channel)
}

function* bootstrap() {
    yield put(updateRoute(history.location))
    history.listen((route) => channel.put(updateRoute(route)))
}

function* dispatchToHistory(
    action: ActionType<typeof push> | ActionType<typeof replace>
): IterableIterator<any> {
    const path = action.payload

    if (isOfType(PUSH_REQUESTED)) {
        history.push(path)
    }

    if (isOfType(REPLACE_REQUESTED)) {
        history.replace(path)
    }
}

//
//   HELPERS and UTILITIES

export const history = process.env.NODE_ENV === 'test' ? createMemoryHistory() : createHashHistory()

export function routeFor(key: string, routeParams = {}) {
    if (!has(ROUTE_MAP, key)) {
        throw new Error(`Invalid route key: ${key}`)
    }
    const params = omitBy(routeParams, isEmpty)
    const route = ROUTE_MAP[key].path
    const url = new URI({ path: compile(route)(params) })
    const keys: Key[] = []

    pathToRegexp(route, keys)
    url.query(omit(params, map(keys, 'name')))

    return url.toString()
}

function matchRoute(location: Location<any>): Route {
    const { pathname, search } = location
    const matched = filter(
        map(ROUTE_MAP, ({ path }, key) => ({
            key,
            ...matchPath(pathname, { path, exact: true })
        })),
        'isExact'
    )

    const routeParams = URI.parseQuery(search) || {}
    const matchedRoute = head(matched)
    const params = {
        ...get(matchedRoute, 'params', {}),
        ...routeParams
    }

    const url = new URI(get(matchedRoute, 'url', pathname))
        .search(params)
        .toString()

    const route = !matchedRoute
        ? {
            url,
            params,
            path: pathname
        }
        : {
            ...matchedRoute,
            url,
            params
        }

    return route as Route
}

export const withRouter = (component: React.ComponentType) => connect(getState)(component)


//
//  SELECTORS

export const getState = (state: RootState) => state.routing
export const getCurrentRoute = getState
export const getParams = createSelector(
    getCurrentRoute,
    (route) => get(route, 'params', {}) as Params
)
