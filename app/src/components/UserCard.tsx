import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { Button } from '@blueprintjs/core'

import { auth, organizations } from '../redux'

import styles from './UserCard.module.scss'

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
            <Button
                text="Logout"
                rightIcon="log-out"
                onClick={ () => dispatch(auth.logout()) }
                small
                minimal
            />
        </div>
    )
}

export default connect()(UserCard)
