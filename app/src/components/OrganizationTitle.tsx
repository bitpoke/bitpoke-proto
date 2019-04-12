import * as React from 'react'
import { connect } from 'react-redux'


import { DispatchProp, api, organizations } from '../redux'

import TitleBar from '../components/TitleBar'
import ResourceActions from '../components/ResourceActions'

type OwnProps = {
    entry?: organizations.IOrganization | null
}

type Props = OwnProps & DispatchProp

const OrganizationTitle: React.SFC<Props> = ({ entry, dispatch }) => {
    const [title, subtitle] = !entry || api.isNewEntry(entry)
        ? ['Create Organization', null]
        : [entry.displayName, entry.name]

    const onDestroy = entry ? () => dispatch(organizations.destroy(entry)) : undefined

    return (
        <TitleBar
            title={ title }
            subtitle={ subtitle }
            actions={
                <ResourceActions
                    entry={ entry }
                    resourceName={ api.Resource.organization }
                    onDestroy={ onDestroy }
                />
            }
        />
    )
}

export default connect()(OrganizationTitle)
