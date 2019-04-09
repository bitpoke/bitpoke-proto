import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { Button, ButtonGroup, Intent } from '@blueprintjs/core'

import { RootState, DispatchProp, api, routing, projects, organizations } from '../redux'

import List from '../components/List'
import TitleBar from '../components/TitleBar'

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
                title={
                    <TitleBar
                        title="Projects"
                        actions={ (
                            <ButtonGroup>
                                <Button
                                    text="Create project"
                                    icon="add"
                                    intent={ Intent.SUCCESS }
                                    onClick={ () =>
                                        dispatch(routing.push(
                                            routing.routeFor('projects', { action: 'new' }))
                                        )
                                    }
                                />
                                <Button
                                    text="Generate random project"
                                    icon="random"
                                    intent={ Intent.SUCCESS }
                                    onClick={ () =>
                                        dispatch(projects.create({
                                            project: {
                                                displayName: faker.commerce.productName()
                                            }
                                        }))
                                    }
                                />
                            </ButtonGroup>
                        ) }
                    />
                }
            />
        </div>
    )
}

export default connect()(ProjectsList)
