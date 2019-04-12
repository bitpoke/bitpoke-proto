import * as React from 'react'
import { connect } from 'react-redux'

import { Navbar as BlueprintNavBar, Spinner, Alignment, Button } from '@blueprintjs/core'


import { RootState, DispatchProp, auth, grpc, routing } from '../redux'

import Link from '../components/Link'
import UserCard from '../components/UserCard'
import OrganizationSwitcher from '../components/OrganizationSwitcher'

import styles from './NavBar.module.scss'

type ReduxProps = {
    currentUser : auth.User,
    isLoading   : boolean
}

type Props = ReduxProps & DispatchProp

const { Group, Heading } = BlueprintNavBar

const NavBar: React.SFC<Props> = (props) => {
    const { currentUser, isLoading, dispatch } = props

    return (
        <BlueprintNavBar>
            <Group align={ Alignment.LEFT }>
                <Heading className={ styles.logo }>
                    <Link to={ routing.routeFor('dashboard') }>Presslabs Dashboard</Link>
                    { isLoading && (
                        <Spinner
                            size={ Spinner.SIZE_SMALL }
                            className={ styles.spinner }
                        />
                    ) }
                </Heading>
                <OrganizationSwitcher />
            </Group>
            <Group align={ Alignment.RIGHT }>
                <UserCard entry={ currentUser } />
                <Button
                    text="Logout"
                    rightIcon="log-out"
                    onClick={ () => dispatch(auth.logout()) }
                    small
                    minimal
                />
            </Group>
        </BlueprintNavBar>
    )
}

const mapStateToProps = (state: RootState): ReduxProps => {
    return {
        currentUser : auth.getCurrentUser(state),
        isLoading   : grpc.isLoading(state)
    }
}

export default connect(mapStateToProps)(NavBar)
