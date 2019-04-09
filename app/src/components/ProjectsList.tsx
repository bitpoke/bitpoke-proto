import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { map } from 'lodash'

import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { RootState, DispatchProp, api, routing, projects } from '../redux'

import Link from '../components/Link'

import styles from './ProjectsList.module.scss'

type ReduxProps = {
    entries: api.ResourcesList<projects.IProject>
}

type Props = ReduxProps & DispatchProp

const ProjectsList: React.SFC<Props> = ({ entries, dispatch }) => {
    return (
        <div>
            <h2>Projects</h2>
            <ButtonGroup>
                <Button
                    text="Create project"
                    icon="add"
                    intent={ Intent.SUCCESS }
                    onClick={ () => dispatch(routing.push(routing.routeFor('project', { action: 'new' }))) }
                />
                <Button
                    text="Create random project"
                    icon="random"
                    intent={ Intent.SUCCESS }
                    onClick={ () => dispatch(projects.create({
                        project: {
                            displayName: faker.commerce.productName()
                        }
                    })) }
                />
            </ButtonGroup>
            <div className={ styles.container }>
                { map(entries, (project) => (
                    <Card
                        key={ project.name }
                        elevation={ Elevation.TWO }
                        interactive
                        onClick={ () => dispatch(routing.push(routing.routeForResource(project))) }
                    >
                        <h5>
                            <Link to={ routing.routeForResource(project) }>{ project.displayName }</Link>
                        </h5>
                        <p>{ project.name }</p>
                    </Card>
                )) }
            </div>
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
