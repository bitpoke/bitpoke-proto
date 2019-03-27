import * as React from 'react'
import { connect } from 'react-redux'
import { createHashHistory, createMemoryHistory, Location } from 'history'
import pathToRegexp, { Key, compile } from 'path-to-regexp'
import { matchPath } from 'react-router'
import { takeEvery, put } from 'redux-saga/effects'
import { channel as createChannel } from 'redux-saga'
import { ActionType, action as createAction } from 'typesafe-actions'
import { createSelector } from 'reselect'

import URI from 'urijs'

import { map, filter, omit, get, head, has, isEmpty } from 'lodash'

import { RootState, app } from '../redux'
import { watchChannel } from '../utils'


//
//  TYPES

export type State = Route | null
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

export const ROUTES = {
    dashboard: {
        path      : '/',
        component : 'DashboardContainer'
    },
    projects: {
        path      : '/projects',
        component : 'ProjectsContainer'
    }
}


//
//  ACTIONS

export const ROUTE_CHANGED = '@ routing / ROUTE_CHANGED'

export const updateRoute = (route: Location<any>) => createAction(ROUTE_CHANGED, route)

const actions = {
    updateRoute
}

//
//  REDUCER

const initialState: State = null

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
    yield watchChannel(channel)
}

function* bootstrap() {
    yield put(updateRoute(history.location))
    history.listen((route) => channel.put(updateRoute(route)))
}

//
//   HELPERS and UTILITIES

export const history = process.env.NODE_ENV === 'test' ? createMemoryHistory() : createHashHistory()
const { push, replace } = history
export { push, replace }

export function routeFor(key: string, params = {}) {
    if (!has(ROUTES, key)) {
        throw new Error(`Invalid route key: ${key}`)
    }

    const route = ROUTES[key].path
    const url = new URI({ path: compile(route)(params) })
    const keys: Key[] = []

    pathToRegexp(route, keys)
    url.query(omit(params, map(keys, 'name')))

    return url.toString()
}

function matchRoute(location: Location<any>): Route {
    const { pathname, search } = location
    const matched = filter(
        map(ROUTES, ({ path }, key) => ({
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
