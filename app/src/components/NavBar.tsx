import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { RootState, auth, projects } from '../redux'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    currentUser: auth.User | null
}

const NavBar: React.SFC<Props> = ({ dispatch }) => {
    return (
        <div>
        </div>
    )
}

const mapStateToProps = (state: RootState): ReduxProps => {
    return {
        currentUser: auth.getCurrentUser(state)
    }
}

export default connect()(NavBar)

