import * as React from 'react'
import { connect } from 'react-redux'
import { Button, Card, Elevation, Intent } from '@blueprintjs/core'

// import Link from '../components/Link'

import { DispatchProp, routing, organizations } from '../redux'

type Props = {
    entry: organizations.IOrganization | null
} & DispatchProp

const OrganizationDetails: React.SFC<Props> = ({ entry, dispatch }) => {
    if (!entry) {
        return null
    }

    return (
        <Card elevation={ Elevation.TWO }>
            <h2>{ entry.displayName }</h2>
            <p>{ entry.name }</p>
            <Button
                text="Edit organization"
                icon="edit"
                intent={ Intent.PRIMARY }
                onClick={ () => dispatch(routing.push(routing.routeForResource(entry, { action: 'edit' }))) }
            />
            <Button
                text="Delete organization"
                icon="trash"
                intent={ Intent.DANGER }
                onClick={ () => dispatch(organizations.destroy(entry)) }
            />
        </Card>
    )
}

export default connect()(OrganizationDetails)
