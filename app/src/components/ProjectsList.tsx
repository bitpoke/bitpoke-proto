import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { Card, Button, ButtonGroup, Intent, Elevation } from '@blueprintjs/core'

import { RootState, DispatchProp, api, routing, projects, organizations } from '../redux'

import Link from '../components/Link'
import List from '../components/List'
import TitleBar from '../components/TitleBar'
import ResourceActions from '../components/ResourceActions'

type OwnProps = {
    organization: organizations.OrganizationName
}

type Props = OwnProps & DispatchProp

const ProjectsList: React.SFC<Props> = (props) => {
    const { organization, dispatch } = props
    return (
        <div>
            <List
                dataRequest={ projects.list({ parent: organization }) }
                dataSelector={ projects.getForOrganization(organization) }
                renderItem={ (entry: projects.IProject) => (
                    <Card
                        key={ entry.name }
                        elevation={ Elevation.TWO }
                        interactive
                        onClick={ () => dispatch(routing.push(routing.routeForResource(entry))) }
                    >
                        <h5>
                            <Link to={ routing.routeForResource(entry) }>{ entry.displayName }</Link>
                        </h5>
                        <p>{ entry.name }</p>
                        <ResourceActions
                            entry={ entry }
                            resourceName={ api.Resource.project }
                            onDestroy={ () => dispatch(projects.destroy(entry)) }
                            withTitles={ false }
                            minimal
                        />
                    </Card>
                ) }
                title={
                    <TitleBar
                        title="Projects"
                        actions={ (
                            <ResourceActions
                                resourceName={ api.Resource.project }
                                onCreate={ () => dispatch(routing.push(
                                    routing.routeFor(api.Resource.project, { slug: '_', action: 'new' })
                                )) }
                                onGenerate={ () => dispatch(projects.create({
                                    project: {
                                        displayName: faker.commerce.productName()
                                    }
                                })) }
                            />
                        ) }
                    />
                }
            />
        </div>
    )
}

export default connect()(ProjectsList)
