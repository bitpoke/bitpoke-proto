import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { map, get } from 'lodash'

import { Card, Button, ButtonGroup, Intent, Elevation } from '@blueprintjs/core'

import { RootState, DispatchProp, api, routing, sites, projects } from '../redux'

import Link from '../components/Link'
import List from '../components/List'
import TitleBar from '../components/TitleBar'
import ResourceActions from '../components/ResourceActions'

type OwnProps = {
    project: projects.ProjectName
}

type Props = OwnProps & DispatchProp

const SitesList: React.SFC<Props> = (props) => {
    const { project, dispatch } = props

    if (!project) {
        return null
    }

    return (
        <List
            dataRequest={ sites.list({ parent: project }) }
            dataSelector={ sites.getForProject(project) }
            renderItem={ (entry: sites.ISite) => (
                <Card
                    key={ entry.name }
                    elevation={ Elevation.TWO }
                    interactive
                    onClick={ () => dispatch(routing.push(routing.routeForResource(entry))) }
                >
                    <h5>
                        <Link to={ routing.routeForResource(entry) }>{ entry.primaryDomain }</Link>
                    </h5>
                    <p>{ entry.name }</p>
                    <ResourceActions
                        entry={ entry }
                        resourceName={ api.Resource.site }
                        onDestroy={ () => dispatch(sites.destroy(entry)) }
                        withTitles={ false }
                        minimal
                    />
                </Card>
            ) }
            title={
                <TitleBar
                    title="Sites"
                    actions={ project && (
                        <ResourceActions
                            resourceName={ api.Resource.site }
                            onCreate={ () =>
                                dispatch(routing.push(
                                    routing.routeFor(api.Resource.site, {
                                        project: projects.parseName(project).slug,
                                        slug: '_',
                                        action: 'new'
                                    })
                                ))
                            }
                            onGenerate={ () =>
                                dispatch(sites.create({
                                    parent: project,
                                    site: {
                                        primaryDomain: faker.internet.domainName()
                                    }
                                }))
                            }
                        />
                    ) }
                />
            }
        />
    )
}

export default connect()(SitesList)
