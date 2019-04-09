import { takeEvery, takeLatest, take, put, call, all, select } from 'redux-saga/effects'
import { ActionType } from 'typesafe-actions'

import { map, compact } from 'lodash'

import { auth, api, routing, organizations, projects, sites, grpc } from '../redux'

export function* saga() {
    yield takeLatest(routing.ROUTE_CHANGED, fetchData)
}

function* fetchData(): IterableIterator<any> {
    const isAuthenticated = yield select(auth.isAuthenticated)
    const organizationMeta = yield select(grpc.getMetadata('organization'))

    if (!isAuthenticated) {
        yield take(auth.ACCESS_GRANTED)
        yield call(fetchData)
        return
    }

    if (!organizationMeta) {
        yield take(organizations.SELECTED)
        yield take(grpc.METADATA_SET)
        yield call(fetchData)
        return
    }

    const route = yield select(routing.getCurrentRoute)

    yield all(compact(map([projects, sites], (scope) => {
        const parsed = scope.parseName(route.url)

        if (parsed && parsed.name) {
            return call(selectOrFetch, scope, parsed.name)
        }

        return null
    })))
}

function* selectOrFetch(scope: any, name: api.ResourceName) {
    const resourceExists = yield select(scope.getByName(name))
    if (resourceExists) {
        return
    }

    yield put(scope.get({ name }))
}
