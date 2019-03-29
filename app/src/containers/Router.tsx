import * as React from 'react'
import { Router as ReactRouter, Switch, Route } from 'react-router-dom'

import { map } from 'lodash'

import { routing } from '../redux'

import * as containers from '../containers'

type Props = {}

const Router: React.SFC<Props> = () => (
    <ReactRouter history={ routing.history }>
        <Switch>
            { map(routing.ROUTE_MAP, ({ path, component }, key) => (
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

export default Router
