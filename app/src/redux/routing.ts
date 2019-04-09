import * as React from 'react'
import { connect } from 'react-redux'
import { createHashHistory, createMemoryHistory, createLocation, Location, Path } from 'history'
import pathToRegexp, { Key, compile } from 'path-to-regexp'
import { matchPath } from 'react-router'
import { channel as createChannel } from 'redux-saga'
import { takeEvery, takeLatest, put, select } from 'redux-saga/effects'
import { ActionType, action as createAction, isOfType } from 'typesafe-actions'
import { createSelector } from 'reselect'

import URI from 'urijs'

import { map, filter, findKey, omit, omitBy, get, head, has, isEmpty, isEqual } from 'lodash'

import { RootState, app, api } from '../redux'
import { watchChannel } from '../utils'


//
//  TYPES

export type State = {
    currentRoute  : Route,
    previousRoute : Route
}

export type Actions = ActionType<typeof actions>
export type Params = object
export type Path = Path

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
    organizations: {
        path      : '/orgs/:slug?/:action?',
        component : 'OrganizationsContainer'
    },
    projects: {
        path      : '/project/:slug?/:action?',
        component : 'ProjectsContainer'
    },
    sites: {
        path      : '/project/:project/site/:slug?/:action?',
        component : 'SitesContainer'
    }
}


//
//  ACTIONS

export const PUSH_REQUESTED    = '@ routing / PUSH_REQUESTED'
export const REPLACE_REQUESTED = '@ routing / REPLACE_REQUESTED'
export const BACK_REQUESTED    = '@ routing / BACK_REQUESTED'
export const PUSH_SKIPPED      = '@ routing / PUSH_SKIPPED'
export const ROUTE_CHANGED     = '@ routing / ROUTE_CHANGED'

export const push = (path: Path) => createAction(PUSH_REQUESTED, path)
export const replace = (path: Path) => createAction(REPLACE_REQUESTED, path)
export const skipPush = (path: Path) => createAction(PUSH_SKIPPED, path)
export const updateRoute = (location: Location<any>) => createAction(ROUTE_CHANGED, location)
export const goBack = () => createAction(BACK_REQUESTED)

const actions = {
    push, replace, updateRoute, goBack
}

//
//  REDUCER

const initialRoute = {
    path   : '/',
    url    : '/',
    key    : findKey(ROUTE_MAP, { path: '/' }),
    params : {}
}

const initialState = {
    currentRoute  : initialRoute,
    previousRoute : initialRoute
}

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case ROUTE_CHANGED: {
            const currentRoute = matchRoute(action.payload)
            const previousRoute = isEqual(currentRoute, state.currentRoute)
                ? state.previousRoute
                : state.currentRoute

            return {
                currentRoute,
                previousRoute
            }
        }

        default:
            return state
    }
}


//
//   SAGA

const channel = createChannel()

export function* saga() {
    yield takeLatest(app.INITIALIZED, bootstrap)
    yield takeEvery([PUSH_REQUESTED, REPLACE_REQUESTED, BACK_REQUESTED], dispatchToHistory)
    yield watchChannel(channel)
}

function* bootstrap() {
    yield put(updateRoute(history.location))
    history.listen((route) => channel.put(updateRoute(route)))
}

function* dispatchToHistory(
    action: ActionType<typeof push> | ActionType<typeof replace> | ActionType<typeof goBack>
): IterableIterator<any> {
    if (isOfType(BACK_REQUESTED, action)) {
        history.goBack()
        return
    }

    const path = action.payload

    if (isOfType(PUSH_REQUESTED, action)) {
        const currentRoute = yield select(getCurrentRoute)
        if (isEqual(currentRoute.url, path)) {
            yield put(skipPush(path))
            return
        }

        history.push(path)
    }

    if (isOfType(REPLACE_REQUESTED, action)) {
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
        .query(omit(params, getPathParams(route)))

    return url.toString()
}

export function routeForResource(resource: api.AnyResourceInstance, params = {}) {
    const pathname = `/${resource.name}`
    const search = URI.buildQuery(params)
    const matchedRoute = matchRoute(createLocation({ pathname, search }))

    if (matchedRoute.key) {
        return routeFor(matchedRoute.key, matchedRoute.params)
    }

    return new URI({ path: pathname, query: search }).toString()
}

function getPathParams(path: Path) {
    const keys: Key[] = []
    pathToRegexp(path, keys)
    return map(keys, 'name')
}

export function matchRoute(location: Location<any>): Route {
    const { pathname, search } = location
    const matched = filter(
        map(ROUTE_MAP, ({ path }, key) => ({
            key,
            ...matchPath(pathname, { path, exact: true })
        })),
        'isExact'
    )

    const routeParams = URI.parseQuery(search) || {}
    const matchedRoute = head(matched) as Route

    const params = {
        ...get(matchedRoute, 'params', {}),
        ...routeParams
    }

    if (matchedRoute) {
        const queryParams = omit(params, getPathParams(matchedRoute.path))
        const url = new URI(matchedRoute.url)
            .query(queryParams)
            .toString()

        const route: Route = {
            ...matchedRoute,
            url,
            params
        }

        return route
    }
    else {
        const url = new URI(pathname)
            .query(params)
            .toString()

        const route: Route = {
            url,
            params,
            path: pathname
        }

        return route
    }
}

export const withRouter = (component: React.ComponentType) => connect(getState)(component)


//
//  SELECTORS

export const getState = (state: RootState) => state.routing
export const getCurrentRoute = createSelector(
    getState,
    (state) => state.currentRoute
)
export const getPreviousRoute = createSelector(
    getState,
    (state) => state.previousRoute
)
export const getParams = createSelector(
    getCurrentRoute,
    (route) => get(route, 'params', {}) as Params
)
