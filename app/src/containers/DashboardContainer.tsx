import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

type Props = {
    dispatch: Dispatch
}

const DashboardContainer: React.SFC<Props> = ({ dispatch }) => {
    return (
        <div></div>
    )
}

export default connect()(DashboardContainer)
