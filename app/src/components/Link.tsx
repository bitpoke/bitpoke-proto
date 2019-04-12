import * as React from 'react'
import { connect } from 'react-redux'
import { isString, isFunction } from 'lodash'

import { DispatchProp, routing } from '../redux'

type Props = {
    to?: string,
    onClick?: () => void,
    className?: string
} & DispatchProp

const Link: React.SFC<Props> = ({ children, to, onClick, className, dispatch }) => {
    const handleClick = isString(to)
        ? (e: React.SyntheticEvent<EventTarget>) => {
            e.preventDefault()
            dispatch(routing.push(to))
        }
        : () => isFunction(onClick) && onClick()

    return (
        <a href={ to } onClick={ handleClick } className={ className }>
            { children }
        </a>
    )
}

export default connect()(Link)
