import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { map, get } from 'lodash'

import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { RootState, DispatchProp, api, routing, sites, projects } from '../redux'

import List from '../components/List'
import TitleBar from '../components/TitleBar'

type OwnProps = {
    project: projects.ProjectName
}

type ReduxProps = {
    parent: projects.IProject | null
}

type Props = OwnProps & ReduxProps & DispatchProp

const SitesList: React.SFC<Props> = (props) => {
    const { project, parent, dispatch } = props

    if (!project) {
        return null
    }

    return (
        <List
            dataRequest={ sites.list({ parent: project }) }
            dataSelector={ sites.getForProject(project) }
            title={
                <TitleBar
                    title="Sites"
                    actions={ parent && (
                        <ButtonGroup>
                            <Button
                                text="Create site"
                                icon="add"
                                intent={ Intent.SUCCESS }
                                onClick={ () =>
                                    dispatch(routing.push(
                                        routing.routeForResource(parent, { action: 'new-site' })
                                    ))
                                }
                            />
                            <Button
                                text="Generate random site"
                                icon="random"
                                intent={ Intent.SUCCESS }
                                onClick={ () =>
                                    dispatch(sites.create({
                                        parent: project,
                                        site: {
                                            primaryDomain: faker.internet.domainName()
                                        }
                                    }))
                                }
                            />
                        </ButtonGroup>
                    ) }
                />
            }
        />
    )
}

function mapStateToProps(state: RootState, ownProps: OwnProps): ReduxProps {
    const parent = projects.getByName(ownProps.project)(state)
    return {
        parent
    }
}

export default connect(mapStateToProps)(SitesList)
