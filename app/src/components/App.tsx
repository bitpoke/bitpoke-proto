import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { RootState, app, auth } from '../redux'

import Router from '../containers/Router'

type Props = {
    dispatch: Dispatch
}

class App extends React.Component<Props> {
    componentDidMount() {
        const { dispatch } = this.props
        dispatch(app.initialize())
    }

    render() {
        return (
            <div>
                <h3>PressLabs Dashboard</h3>
                <Router />
            </div>
        )
    }
}

export default connect()(App)
