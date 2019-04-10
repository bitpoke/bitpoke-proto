import * as React from 'react'
import { connect } from 'react-redux'
import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState, DispatchProp, api, routing, sites } from '../redux'

import TitleBar from '../components/TitleBar'
import SitesList from '../components/SitesList'
import ResourceActions from '../components/ResourceActions'

type OwnProps = {
    entry?: sites.ISite | null
}

type Props = OwnProps & DispatchProp

const SiteTitle: React.SFC<Props> = ({ entry, dispatch }) => {
    const [title, subtitle] = !entry || api.isNewEntry(entry)
        ? ['Create Site', null]
        : [entry.primaryDomain, entry.name]

    const onDestroy = entry ? () => dispatch(sites.destroy(entry)) : undefined

    return (
        <TitleBar
            title={ title }
            subtitle={ subtitle }
            actions={
                <ResourceActions
                    entry={ entry }
                    resourceName={ api.Resource.site }
                    onDestroy={ onDestroy }
                />
            }
        />
    )
}

export default connect()(SiteTitle)
