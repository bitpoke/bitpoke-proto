import * as React from 'react'
import { connect } from 'react-redux'

import { DispatchProp, api, routing, projects } from '../redux'

import TitleBar from '../components/TitleBar'
import ResourceActions from '../components/ResourceActions'

type OwnProps = {
    entry?: projects.IProject | null,
    title?: string | null,
    withActionTitles?: boolean,
    withMinimalActions?: boolean
}

type Props = OwnProps & DispatchProp

const ProjectTitle: React.SFC<Props> = (props) => {
    const { entry, withActionTitles, withMinimalActions, dispatch } = props
    const [title, subtitle, link, onDestroy] = !entry || api.isNewEntry(entry)
        ? [props.title || 'Create Project', null, null, undefined]
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
                    withTitles={ withActionTitles }
                    minimal={ withMinimalActions }
                />
            }
        />
    )
}

ProjectTitle.defaultProps = {
    withActionTitles: true,
    withMinimalActions: false
}

export default connect()(ProjectTitle)
