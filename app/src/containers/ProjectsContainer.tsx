import * as React from 'react'
import { connect } from 'react-redux'
import { Switch, Route } from 'react-router-dom'

import Container from '../components/Container'
import ProjectForm from '../components/ProjectForm'
import ProjectDetails from '../components/ProjectDetails'

import { RootState, routing, projects } from '../redux'

type ReduxProps = {
    project: projects.IProject | null
}

type Props = ReduxProps

const ProjectsContainer: React.SFC<Props> = ({ project }) => {
    return (
        <Container>
            <Switch>
                <Route
                    path={ routing.routeFor('project', { action: 'new' }) }
                    component={ ProjectForm }
                />
                { project && (
                    <Route
                        path={ routing.routeForResource(project, { action: 'edit' }) }
                        render={ () => <ProjectForm initialValues={ { project } } /> }
                    />
                ) }
                <Route
                    path={ routing.routeFor('project') }
                    render={ () => <ProjectDetails entry={ project } /> }
                />
            </Switch>
        </Container>
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    return {
        project: projects.getForCurrentURL(state)
    }
}

export default connect(mapStateToProps)(ProjectsContainer)
