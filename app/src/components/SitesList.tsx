import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'


import { Card, Elevation } from '@blueprintjs/core'

import { DispatchProp, api, routing, sites, projects } from '../redux'

import List from '../components/List'
import TitleBar from '../components/TitleBar'
import SiteTitle from '../components/SiteTitle'
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
                    <SiteTitle entry={ entry } withActionTitles={ false } withMinimalActions />
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
