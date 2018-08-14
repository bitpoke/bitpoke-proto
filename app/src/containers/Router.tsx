import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'
import { Router as ReactRouter, Switch, Route, Redirect } from 'react-router-dom'

import { map } from 'lodash'

import { RootState, routing } from '../redux'

import * as containers from '../containers'

type Props = {
    dispatch: Dispatch
}

class Router extends React.Component<Props> {
    componentDidMount() {
        const { dispatch } = this.props
        const { history, updateRoute } = routing
        history.listen(({ pathname }) => dispatch(updateRoute(pathname)))
        dispatch(updateRoute(history.location.pathname))
    }

    render() {
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

export default connect()(Router)
