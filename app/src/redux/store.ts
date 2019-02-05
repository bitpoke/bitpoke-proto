import { createStore as createReduxStore, applyMiddleware, compose, combineReducers } from 'redux'
import { createLogger } from 'redux-logger'
import { persistStore, persistReducer } from 'redux-persist'
import localForage from 'localforage'
import createSagaMiddleware, { SagaIterator } from 'redux-saga'
import { fork, all } from 'redux-saga/effects'

import { map, reduce, compact, has } from 'lodash'

import * as reduxModules from '../redux'

const reducers = reduce(reduxModules, (acc, module, key) => {
    if (!has(module, 'reducer')) {
        return acc
    }
    return {
        ...acc,
        [key]: module.reducer
    }
}, {})

const rootReducer = persistReducer({
    key       : 'root',
    whitelist : ['auth'],
    storage   : localForage
}, combineReducers(reducers))

const initialState = rootReducer(undefined, {} as any)
const sagas = compact(map(reduxModules, 'saga'))

function* rootSaga(): SagaIterator {
    yield all(map(sagas, (saga) => fork(saga)))
}

function createStore() {
    const middleware = []
    const enhancers = [] as any

    if (process.env.NODE_ENV === 'development') {
        const logger = createLogger({
            level     : 'info',
            collapsed : true,
            timestamp : false,
            duration  : true
        })
        middleware.push(logger)

        // const w : any = window as any
        // const devTools: any = w.devToolsExtension ? w.devToolsExtension() : (f:any) => f
        // if (isFunction(devTools)) {
        //     enhancers.push(devTools())
        // }
    }

    const sagaMiddleware = createSagaMiddleware()
    middleware.push(sagaMiddleware)

    const store = createReduxStore(
        rootReducer,
        initialState,
        compose(
            applyMiddleware(...middleware),
            ...enhancers
        )
    )

    sagaMiddleware.run(rootSaga)

    return store
}

export default () => {
    const store = createStore()
    const persistor = persistStore(store)
    return { store, persistor }
}
