import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { RootState, app, auth } from '../redux'

import Router from '../containers/Router'
import NavBar from '../components/NavBar'

import './App.css'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    isAuthenticated: boolean
}

class App extends React.Component<Props & ReduxProps> {
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
            <div className="App_container">
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
