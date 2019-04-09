import * as React from 'react'
import { connect } from 'react-redux'

import { RootState, DispatchProp, app, auth } from '../redux'

import Router from '../containers/Router'
import NavBar from '../components/NavBar'

import './App.scss'

type ReduxProps = {
    isAuthenticated: boolean
}

type Props = ReduxProps & DispatchProp

class App extends React.Component<Props> {
    componentDidMount() {
        const { dispatch } = this.props
        dispatch(app.initialize())
    }

    render() {
        const { isAuthenticated } = this.props
        if (!isAuthenticated) {
            return null
        }

        return (
            <div>
                <NavBar />
                <Router />
            </div>
        )
    }
}

const mapStateToProps = (state: RootState): ReduxProps => {
    return {
        isAuthenticated: auth.isAuthenticated(state)
    }
}

export default connect(mapStateToProps)(App)
