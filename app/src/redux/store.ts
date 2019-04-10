import localforage from 'localforage'
import Raven from 'raven-js'
import { Middleware } from 'redux'
import createFilter from 'redux-persist-transform-filter'

import * as storeUtils from '../utils/store'
import * as modules from '../redux'

const persistedReducers = [
    'auth',
    'organizations'
]

const authFilter = createFilter('auth', ['token'])
const orgsFilter = createFilter('organizations', ['current'])

const persistConfig = {
    key        : 'root',
    whitelist  : persistedReducers,
    storage    : localforage,
    transforms : [authFilter, orgsFilter]
}

const middleware: Array<Middleware<any>> = []

if (process.env.NODE_ENV !== 'development' && process.env.REACT_APP_SENTRY_DSN) {
    Raven.config(process.env.REACT_APP_SENTRY_DSN).install()
    const sentryMiddleware: Middleware = storeUtils.createSentryMiddleware(Raven)
    middleware.push(sentryMiddleware)
}

const initialState = {}
const rootReducer = storeUtils.createRootReducer(modules, persistConfig)
const rootSaga = storeUtils.createRootSaga(modules)

export const { store, persistor } = storeUtils.createPersistedStore(
    rootReducer,
    rootSaga,
    initialState,
    middleware
)
