import * as React from 'react'
import { connect } from 'react-redux'
import { Button, Card, Elevation, Intent } from '@blueprintjs/core'

import { DispatchProp, routing, projects } from '../redux'

import SitesList from '../components/SitesList'

type Props = {
    entry: projects.IProject | null
} & DispatchProp

const ProjectDetails: React.SFC<Props> = ({ entry, dispatch }) => {
    if (!entry) {
        return null
    }

    return (
        <div>
            <Card elevation={ Elevation.TWO }>
                <h2>{ entry.displayName }</h2>
                <p>{ entry.name }</p>
                <Button
                    text="Edit project"
                    icon="edit"
                    intent={ Intent.PRIMARY }
                    onClick={ () => dispatch(routing.push(routing.routeForResource(entry, { action: 'edit' }))) }
                />
                <Button
                    text="Delete project"
                    icon="trash"
                    intent={ Intent.DANGER }
                    onClick={ () => dispatch(projects.destroy(entry)) }
                />
            </Card>
            <SitesList project={ entry.name } />
        </div>
    )
}

export default connect()(ProjectDetails)
