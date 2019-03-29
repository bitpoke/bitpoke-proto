import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'
import { isString, isFunction } from 'lodash'

import { routing } from '../redux'

type Props = {
    dispatch: Dispatch,
    to?: string,
    onClick?: () => void
}

const Link: React.SFC<Props> = ({ children, to, onClick, dispatch }) => {
    const handleClick = isString(to)
        ? (e: React.SyntheticEvent<EventTarget>) => {
            e.preventDefault()
            dispatch(routing.push(to))
        }
        : () => isFunction(onClick) && onClick()

    return (
        <a href={ to } onClick={ handleClick }>
            { children }
        </a>
    )
}

export default connect()(Link)
