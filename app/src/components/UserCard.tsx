import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { auth } from '../redux'

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
            <button onClick={ () => dispatch(auth.refreshToken()) }>Refresh Token</button>
            <button onClick={ () => dispatch(auth.logout()) }>Logout</button>
        </div>
    )
}

export default connect()(UserCard)
