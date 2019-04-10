import * as React from 'react'
import { connect } from 'react-redux'
import { Button, Card, Elevation, Intent } from '@blueprintjs/core'

import OrganizationTitle from '../components/OrganizationTitle'

import { DispatchProp, routing, organizations } from '../redux'

type Props = {
    entry: organizations.IOrganization | null
} & DispatchProp

const OrganizationDetails: React.SFC<Props> = (props) => {
    const { entry, dispatch } = props

    if (!entry) {
        return null
    }

    return (
        <div>
            <OrganizationTitle entry={ entry } />
        </div>
    )
}

export default connect()(OrganizationDetails)
