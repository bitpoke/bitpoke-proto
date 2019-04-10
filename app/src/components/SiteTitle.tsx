import * as React from 'react'
import { connect } from 'react-redux'


import { DispatchProp, api, routing, sites } from '../redux'

import TitleBar from '../components/TitleBar'
import ResourceActions from '../components/ResourceActions'
import SiteStatusTag from '../components/SiteStatusTag'

type OwnProps = {
    entry?: sites.ISite | null,
    title?: string | null,
    withActionTitles?: boolean,
    withMinimalActions?: boolean
}

type Props = OwnProps & DispatchProp

const SiteTitle: React.SFC<Props> = (props) => {
    const { entry, withActionTitles, withMinimalActions, dispatch } = props
    const [title, subtitle, link, onDestroy] = !entry || api.isNewEntry(entry)
        ? [props.title || 'Create Site', null, null, undefined]
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
                    withTitles={ withActionTitles }
                    minimal={ withMinimalActions }
                />
            }
        />
    )
}

SiteTitle.defaultProps = {
    withActionTitles: true,
    withMinimalActions: false
}

export default connect()(SiteTitle)
