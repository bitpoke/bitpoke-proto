import * as React from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'
import faker from 'faker'

import { map } from 'lodash'

import { Button, Card, Elevation, Intent } from '@blueprintjs/core'

import { RootState, api, routing, projects } from '../redux'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    entries: api.ResourcesList<projects.IProject>
}

const ProjectsList: React.SFC<Props & ReduxProps> = ({ entries, dispatch }) => {
    return (
        <div>
            <h2>Projects</h2>
            <Button
                text="Create project"
                icon="add"
                intent={ Intent.SUCCESS }
                onClick={ () => routing.push(routing.routeFor('onboarding', { step: 'project' })) }
            />
            <Button
                text="Create random project"
                icon="random"
                intent={ Intent.SUCCESS }
                onClick={ () => dispatch(projects.create({
                    displayName: faker.commerce.productName()
                })) }
            />
            { map(entries, (project) => (
                <Card
                    key={ `project-${project.name}` }
                    elevation={ Elevation.TWO }
                >
                    <h5><a href="#">{ project.displayName }</a></h5>
                    <p>{ project.name }</p>
                    <Button
                        text="Delete project"
                        icon="trash"
                        intent={ Intent.DANGER }
                        onClick={ () => dispatch(projects.destroy(project)) }
                    />
                </Card>
            )) }
        </div>
    )
}

function mapStateToProps(state: RootState) {
    const entries = projects.getForCurrentOrganization(state)
    return {
        entries
    }
}

export default connect(mapStateToProps)(ProjectsList)
