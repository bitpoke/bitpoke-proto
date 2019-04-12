import * as React from 'react'

import OrganizationTitle from '../components/OrganizationTitle'

import { organizations } from '../redux'

type Props = {
    entry: organizations.IOrganization | null
}

const OrganizationDetails: React.SFC<Props> = (props) => {
    const { entry } = props

    if (!entry) {
        return null
    }

    return (
        <div>
            <OrganizationTitle entry={ entry } />
        </div>
    )
}

export default OrganizationDetails
