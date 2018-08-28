import * as React from 'react'
import * as ReactDOM from 'react-dom'

import { Provider } from 'react-redux'
import { PersistGate } from 'redux-persist/integration/react'

import createStore from './redux/store'
import registerServiceWorker from './registerServiceWorker'

import App from './components/App'

import './index.css'

const { store, persistor } = createStore()

ReactDOM.render(
    <PersistGate persistor={ persistor } loading={ <div /> }>
        <Provider store={ store }>
            <App />
        </Provider>
    </PersistGate>,
    document.getElementById('root') as HTMLElement
)

registerServiceWorker()

    // <PersistGate persistor={ persistor } loading={ <div /> }>
    //     <Provider store={ store }>
    //         <App />
    //     </Provider>
    // </PersistGate>,
