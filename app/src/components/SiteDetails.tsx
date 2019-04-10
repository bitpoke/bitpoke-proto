import * as React from 'react'

import { projects } from '../redux'

import SiteTitle from '../components/SiteTitle'

type Props = {
    entry: projects.IProject | null
}

const SiteDetails: React.SFC<Props> = ({ entry }) => {
    if (!entry) {
        return null
    }

    return (
        <div>
            <SiteTitle entry={ entry } />
        </div>
    )
}

export default SiteDetails
