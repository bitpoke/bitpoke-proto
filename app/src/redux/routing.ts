import * as React from 'react'
import { connect } from 'react-redux'
import { createHashHistory, createMemoryHistory } from 'history'
import pathToRegexp, { Key, compile } from 'path-to-regexp'
import { matchPath } from 'react-router'
import { takeEvery, put, select, take, call } from 'redux-saga/effects'
import { createSelector } from 'reselect'
import { ActionType, action as createAction } from 'typesafe-actions'

import URI from 'urijs'

import { map, filter, omit, head, has, isEmpty, startCase } from 'lodash'

import { RootState } from '../redux'


//
//  TYPES

export type State = Route
export type Actions = ActionType<typeof actions>

export type Route = {
    key: string
} | null


//
//  ROUTES

export const ROUTES = {
    dashboard: {
        path      : '/',
        component : 'DashboardContainer'
    }
}


//
//  ACTIONS

export const ROUTE_CHANGED = '@ routing / ROUTE_CHANGED'

export const updateRoute = (route: string) => createAction(ROUTE_CHANGED, route)

const actions = {
    updateRoute
}

//
//  REDUCER

export const initialState: State = null

export function reducer(state: State = initialState, action: Actions) {
    switch (action.type) {
        case ROUTE_CHANGED: {
            return matchRoute(action.payload) || {}
        }

        default:
            return state
    }
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

function matchRoute(pathname: string) {
    const matched = filter(
        map(ROUTES, ({ path }, key) => ({
            key,
            ...matchPath(pathname, { path, exact: true })
        })),
        'isExact'
    )

    return head(matched)
}

export const withRouter = (component: React.ComponentType) => connect(getState)(component)


//
//  SELECTORS

export const getState = (state: RootState) => (state.routing)
