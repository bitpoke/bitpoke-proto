import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { auth, routing } from '../redux'

import './UserCard.css'

type Props = {
    dispatch: Dispatch,
    entry: auth.User
}

const UserCard: React.SFC<Props> = ({ entry, dispatch }) => {
    return (
        <div className="userCard_container">
            <img className="userCard_avatar" src={ entry.avatarURL } />
            <strong className="userCard_email">{ entry.email }</strong>
            <button onClick={ () => dispatch(auth.refreshToken()) }>Refresh Token</button>
            <button onClick={ () => dispatch(auth.logout()) }>Logout</button>
        </div>
    )
}

export default connect()(UserCard)
