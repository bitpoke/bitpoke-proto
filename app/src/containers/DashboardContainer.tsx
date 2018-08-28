import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { projects } from '../redux'

type Props = {
    dispatch: Dispatch
}

const DashboardContainer: React.SFC<Props> = ({ dispatch }) => {
    return (
        <div>
            <h4>Projects</h4>
            <button onClick={ () => dispatch(projects.list()) }>
                List Projects
            </button>
        </div>
    )
}

export default connect()(DashboardContainer)
