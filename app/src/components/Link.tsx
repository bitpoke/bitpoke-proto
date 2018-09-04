import * as React from 'react'

import { isString, isFunction } from 'lodash'

import { routing } from '../redux'

type Props = {
    to?: string
    onClick?: () => void
}

const Link: React.SFC<Props> = ({ children, to, onClick }) => {
    const handleClick = isString(to)
        ? (e: React.SyntheticEvent<EventTarget>) => { e.preventDefault() && routing.push(to) }
        : () => isFunction(onClick) && onClick()

    return (
        <a href={ to } onClick={ handleClick }>
            { children }
        </a>
    )
}

export default Link
