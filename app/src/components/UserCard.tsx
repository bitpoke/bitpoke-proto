import * as React from 'react'
import { connect } from 'react-redux'

import { RootState, DispatchProp, auth, routing, organizations } from '../redux'

import Link from '../components/Link'

import styles from './UserCard.module.scss'

type ReduxProps = {
    currentOrganization: organizations.IOrganization | null
}

type OwnProps = {
    entry: auth.User
}

type Props = OwnProps & ReduxProps & DispatchProp

const UserCard: React.SFC<Props> = ({ entry, currentOrganization, dispatch }) => {
    if (!entry) {
        return null
    }

    return (
        <Link
            to={ currentOrganization ? routing.routeForResource(currentOrganization) : routing.routeFor('dashboard') }
            className={ styles.container }
        >
            <img className={ styles.avatar } src={ entry.avatarURL } />
            <strong className={ styles.email }>{ entry.email }</strong>
        </Link>
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    const currentOrganization = organizations.getCurrent(state)
    return {
        currentOrganization
    }
}

export default connect(mapStateToProps)(UserCard)
