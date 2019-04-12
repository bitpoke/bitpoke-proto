import * as React from 'react'

import { projects } from '../redux'

import ProjectTitle from '../components/ProjectTitle'
import SitesList from '../components/SitesList'

type Props = {
    entry: projects.IProject | null
}

const ProjectDetails: React.SFC<Props> = ({ entry }) => {
    if (!entry) {
        return null
    }

    return (
        <div>
            <ProjectTitle entry={ entry } />
            <SitesList project={ entry.name } />
        </div>
    )
}

export default ProjectDetails
