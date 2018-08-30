import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import ProjectsList from '../components/ProjectsList'

type Props = {
    dispatch: Dispatch
}

const DashboardContainer: React.SFC<Props> = ({ dispatch }) => {
    return (
        <div>
            <ProjectsList />
        </div>
    )
}

export default connect()(DashboardContainer)
