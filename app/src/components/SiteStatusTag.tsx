import * as React from 'react'
import { Tag, Tooltip, Intent, Position } from '@blueprintjs/core'

import { get } from 'lodash'

import { sites } from '../redux'

type OwnProps = {
    entry?: sites.ISite | null
}

type Props = OwnProps

const STATUS_INTENT = {
    [sites.Status.UNSPECIFIED]  : Intent.NONE,
    [sites.Status.PROVISIONING] : Intent.WARNING,
    [sites.Status.RUNNING]      : Intent.PRIMARY,
    [sites.Status.ERROR]        : Intent.DANGER
}

const SiteStatusTag: React.SFC<Props> = ({ entry }) => {
    if (!entry) {
        return null
    }

    const statusName = sites.statusName(entry.status)
    const statusMessage = entry.statusMessage || `Current site status: ${statusName}`
    const statusIntent = get(STATUS_INTENT, entry.status || '', Intent.NONE)

    return (
        <Tooltip content={ statusMessage } position={ Position.RIGHT }>
            <Tag minimal intent={ statusIntent }>{ statusName }</Tag>
        </Tooltip>
    )
}

export default SiteStatusTag
