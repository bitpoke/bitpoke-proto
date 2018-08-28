import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'
import { Router as ReactRouter, Switch, Route, Redirect } from 'react-router-dom'

import { map } from 'lodash'

import { RootState, auth, routing } from '../redux'

import * as containers from '../containers'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    isAuthenticated: boolean
}

class Router extends React.Component<Props & ReduxProps> {
    componentDidMount() {
        const { dispatch } = this.props
        const { history, updateRoute } = routing
        history.listen(({ pathname }) => dispatch(updateRoute(pathname)))
        dispatch(updateRoute(history.location.pathname))
    }

    render() {
        const { isAuthenticated } = this.props
        if (!isAuthenticated) {
            return null
        }

        return (
            <ReactRouter history={ routing.history }>
                <Switch>
                    { map(routing.ROUTES, ({ path, component }, key) => (
                        <Route
                            key={ key }
                            path={ path }
                            component={ containers[component] }
                            exact
                        />
                    )) }
                </Switch>
            </ReactRouter>
        )
    }
}

const mapStateToProps = (state: RootState): ReduxProps => {
    return {
        isAuthenticated: auth.isAuthenticated(state)
    }
}

export default connect(mapStateToProps)(Router)
