import { put, take } from 'redux-saga/effects'
import { Channel } from 'redux-saga'

export function* watchChannel(actionChannel: Channel<{}>) {
    while (true) {
        yield put(yield take(actionChannel))
    }
}
