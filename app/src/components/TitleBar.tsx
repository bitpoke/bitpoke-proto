import * as React from 'react'

import { isFunction } from 'lodash'

import { Card, Elevation, Button, Intent, Tag, Classes } from '@blueprintjs/core'

import styles from './TitleBar.module.scss'

export type Props = {
    title?: string | null,
    subtitle?: string | null,
    tag?: React.ReactNode,
    actions?: React.ReactNode
}

const TitleBar: React.SFC<Props> = (props) => {
    const { title, subtitle, tag, actions } = props

    return (
        <div className={ styles.container }>
            <div className={ styles.title }>
                <h2>
                    { title }
                    { tag && <Tag minimal round>{ tag }</Tag> }
                </h2>
                { subtitle && <h4 className={ Classes.TEXT_MUTED }>{ subtitle }</h4> }
            </div>
            <div className={ styles.actions }>
                { actions }
            </div>
        </div>
    )
}

export default TitleBar
