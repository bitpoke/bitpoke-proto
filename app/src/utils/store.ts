import {
    createStore as createReduxStore, combineReducers, applyMiddleware, compose,
    Reducer
} from 'redux'
import { createLogger } from 'redux-logger'

import {
    persistStore, purgeStoredState, persistCombineReducers,
    PersistConfig, PersistedState
} from 'redux-persist'

import createSagaMiddleware, { Saga } from 'redux-saga'
import { RavenStatic } from 'raven-js'
import createRavenMiddleware from 'raven-for-redux'

import { fork, all } from 'redux-saga/effects'

import { map, reduce, pickBy, get, compact, noop, has, isEmpty, isFunction } from 'lodash'

import { auth, RootState, AnyAction } from '../redux'

export function createSentryMiddleware(sentryClient: RavenStatic) {
    return createRavenMiddleware(sentryClient, {
        getUserContext: (state) => {
            const user = auth.getCurrentUser(state)
            if (!user) {
                return {}
            }

            const { id, email, name } = user

            return {
                id,
                email,
                username: name
            }
        }
    })
}

export function createRootReducer(modules: object, persistConfig: PersistConfig) {
    const reducers = reduce(modules, (acc, module, key) => {
        if (!isFunction(get(module, 'reducer'))) {
            return acc
        }
        return {
            ...acc,
            [key]: get(module, 'reducer')
        }
    }, {})

    if (isEmpty(persistConfig)) {
        return combineReducers(reducers)
    }

    const rootReducer: Reducer = persistCombineReducers(persistConfig, reducers)

    return (state: RootState | undefined, action: AnyAction): RootState => {
        if (action.type === auth.LOGOUT_REQUESTED) {
            purgeStoredState(persistConfig)
            return rootReducer(undefined, action)
        }

        return rootReducer(state, action)
    }
}

export function createRootSaga(modules: object) {
    const sagas = map(modules, (m) => get(m, 'saga', noop))

    return function* rootSaga() {
        yield all(map(sagas, (saga) => fork(saga)))
    }
}

export function createStore(
    rootReducer: Reducer,
    rootSaga: Saga<any[]>,
    initialState = {},
    customMiddleware = []
) {
    const middleware = []

    if (process.env.NODE_ENV === 'development') {
        const logger = createLogger({
            level     : 'info',
            collapsed : true,
            timestamp : false,
            duration  : true
        })
        middleware.push(logger)
    }

    const sagaMiddleware = createSagaMiddleware()
    middleware.push(sagaMiddleware)

    const store = createReduxStore(
        rootReducer,
        initialState,
        compose(
            applyMiddleware(...middleware, ...customMiddleware)
        )
    )

    sagaMiddleware.run(rootSaga)

    return store
}

export function createPersistedStore(
    rootReducer: any,
    rootSaga: Saga<any[]>,
    initialState = {},
    customMiddleware = []
) {
    const store = createStore(rootReducer, rootSaga, initialState, customMiddleware)
    const persistor = persistStore(store)
    return { store, persistor }
}
