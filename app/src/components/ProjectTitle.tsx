import * as React from 'react'
import { connect } from 'react-redux'
import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState, DispatchProp, api, routing, projects } from '../redux'

import TitleBar from '../components/TitleBar'
import SitesList from '../components/SitesList'
import ResourceActions from '../components/ResourceActions'

type OwnProps = {
    entry?: projects.IProject | null
}

type Props = OwnProps & DispatchProp

const ProjectTitle: React.SFC<Props> = ({ entry, dispatch }) => {
    const [title, subtitle, link, onDestroy] = !entry || api.isNewEntry(entry)
        ? ['Create Project', null, null, undefined]
        : [entry.displayName, entry.name, routing.routeForResource(entry), () => dispatch(projects.destroy(entry))]

    return (
        <TitleBar
            title={ title }
            subtitle={ subtitle }
            link={ link }
            actions={
                <ResourceActions
                    entry={ entry }
                    resourceName={ api.Resource.project }
                    onDestroy={ onDestroy }
                />
            }
        />
    )
}

export default connect()(ProjectTitle)

