import * as React from 'react'
import { connect } from 'react-redux'
import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState, DispatchProp, api, routing, sites } from '../redux'

import TitleBar from '../components/TitleBar'
import SitesList from '../components/SitesList'
import ResourceActions from '../components/ResourceActions'
import SiteStatusTag from '../components/SiteStatusTag'

type OwnProps = {
    entry?: sites.ISite | null
}

type Props = OwnProps & DispatchProp

const SiteTitle: React.SFC<Props> = ({ entry, dispatch }) => {
    const [title, subtitle, link, onDestroy] = !entry || api.isNewEntry(entry)
        ? ['Create Site', null, null, undefined]
        : [entry.primaryDomain, entry.name, routing.routeForResource(entry), () => dispatch(sites.destroy(entry))]

    return (
        <TitleBar
            title={ title }
            subtitle={ subtitle }
            link={ link }
            tag={ <SiteStatusTag entry={ entry } /> }
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
