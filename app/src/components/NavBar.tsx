import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { RootState, auth, projects, routing } from '../redux'

import Link from '../components/Link'
import UserCard from '../components/UserCard'

import './NavBar.css'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    currentUser: auth.User
}

const NavBar: React.SFC<Props & ReduxProps> = ({ dispatch, currentUser }) => {
    return (
        <div className="NavBar_container">
            <h2 className="NavBar_logo">
                <Link to={ routing.routeFor('dashboard') }>Dashboard</Link>
            </h2>
            <UserCard entry={ currentUser } />
        </div>
    )
}

const mapStateToProps = (state: RootState): ReduxProps => {
    return {
        currentUser: auth.getCurrentUser(state)
    }
}

export default connect(mapStateToProps)(NavBar)
