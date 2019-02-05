import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { auth, organizations } from '../redux'

import styles from './UserCard.module.scss'

import {
    Organization,
    CreateOrganizationRequest
} from '@presslabs/dashboard-proto'

type Props = {
    dispatch: Dispatch,
    entry: auth.User
}

const UserCard: React.SFC<Props> = ({ entry, dispatch }) => {
    if (!entry) {
        return null
    }

    return (
        <div className={ styles.container }>
            <img className={ styles.avatar } src={ entry.avatarURL } />
            <strong className={ styles.email }>{ entry.email }</strong>
            <button onClick={ () => dispatch(auth.refreshToken()) }>Refresh Token</button>
            <button onClick={ () => dispatch(organizations.list()) }>List orgs</button>
            <button onClick={ () => {
                const organization = new Organization()
                organization.setDisplayName('A New Organization')
                const request = new CreateOrganizationRequest()
                request.setOrganization(organization)
                dispatch(organizations.create(request))
            } }>Create org</button>
            <button onClick={ () => dispatch(auth.logout()) }>Logout</button>
        </div>
    )
}

export default connect()(UserCard)
