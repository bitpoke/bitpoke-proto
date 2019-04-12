import * as React from 'react'
import * as ReactDOM from 'react-dom'

import { Provider } from 'react-redux'
import { PersistGate } from 'redux-persist/integration/react'

import { store, persistor } from './redux/store'
import registerServiceWorker from './registerServiceWorker'

import App from './components/App'

import './index.scss'

ReactDOM.render((
    <PersistGate persistor={ persistor } loading={ <div /> }>
        <Provider store={ store }>
            <App />
        </Provider>
    </PersistGate>
), document.getElementById('root') as HTMLElement)

registerServiceWorker()
