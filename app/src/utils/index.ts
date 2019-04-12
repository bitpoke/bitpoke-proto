import { put, take } from 'redux-saga/effects'
import { Channel } from 'redux-saga'

export type Omit<T, K> = Pick<T, Exclude<keyof T, K>>

export function* watchChannel(actionChannel: Channel<{}>) {
    while (true) {
        yield put(yield take(actionChannel))
    }
}
