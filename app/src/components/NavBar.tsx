import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { Navbar as BlueprintNavBar, Alignment } from '@blueprintjs/core'

import { RootState, auth, routing } from '../redux'

import Link from '../components/Link'
import UserCard from '../components/UserCard'

import styles from './NavBar.module.scss'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    currentUser: auth.User
}

const { Group, Heading } = BlueprintNavBar

const NavBar: React.SFC<Props & ReduxProps> = ({ dispatch, currentUser }) => {
    return (
        <BlueprintNavBar>
            <Group align={ Alignment.LEFT }>
                <Heading>
                    <Link to={ routing.routeFor('dashboard') }></Link>
                </Heading>
            </Group>
            <Group align={ Alignment.RIGHT }>
                <UserCard entry={ currentUser } />
            </Group>
        </BlueprintNavBar>
    )
}

const mapStateToProps = (state: RootState): ReduxProps => {
    return {
        currentUser: auth.getCurrentUser(state)
    }
}

export default connect(mapStateToProps)(NavBar)
