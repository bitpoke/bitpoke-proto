import * as React from 'react'
import { connect } from 'react-redux'

import { Navbar as BlueprintNavBar, Spinner, Alignment, Button } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState, DispatchProp, auth, grpc, routing, organizations } from '../redux'

import Link from '../components/Link'
import UserCard from '../components/UserCard'
import OrganizationsList from '../components/OrganizationsList'

import styles from './NavBar.module.scss'

type ReduxProps = {
    currentUser         : auth.User,
    currentOrganization : organizations.IOrganization | null,
    isLoading           : boolean
}

type Props = ReduxProps & DispatchProp

const { Group, Heading } = BlueprintNavBar

const NavBar: React.SFC<Props> = (props) => {
    const { currentUser, currentOrganization, isLoading, dispatch } = props
    return (
        <BlueprintNavBar>
            <Group align={ Alignment.LEFT }>
                <Heading className={ styles.logo }>
                    <Link to={ routing.routeFor('dashboard', {
                        org: get(currentOrganization, 'name', null)
                    }) }>
                        Presslabs Dashboard
                    </Link>
                    { isLoading && (
                        <Spinner
                            size={ Spinner.SIZE_SMALL }
                            className={ styles.spinner }
                        />
                    ) }
                </Heading>
                <OrganizationsList />
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
        currentUser         : auth.getCurrentUser(state),
        currentOrganization : organizations.getCurrent(state),
        isLoading           : grpc.isLoading(state)
    }
}

export default connect(mapStateToProps)(NavBar)
