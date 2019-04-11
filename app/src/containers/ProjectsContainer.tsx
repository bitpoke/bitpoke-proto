import * as React from 'react'
import { connect } from 'react-redux'
import { Switch, Route } from 'react-router-dom'

import Container from '../components/Container'
import ProjectForm from '../components/ProjectForm'
import ProjectDetails from '../components/ProjectDetails'

import { RootState, DispatchProp, routing, organizations, projects } from '../redux'

type ReduxProps = {
    project: projects.IProject | null
}

type Props = ReduxProps & DispatchProp

const ProjectsContainer: React.SFC<Props> = (props) => {
    const { project } = props

    return (
        <Container>
            <Switch>
                <Route
                    path={ routing.routeFor('projects', { slug: '_', action: 'new' }) }
                    render={ () => <ProjectForm initialValues={ {} } /> }
                />
                { project && (
                    <Route
                        path={ routing.routeForResource(project, { action: 'edit' }) }
                        render={ () => <ProjectForm initialValues={ { project } } /> }
                    />
                ) }
                { project && (
                    <Route
                        path={ routing.routeForResource(project) }
                        render={ () => <ProjectDetails entry={ project } /> }
                    />
                ) }
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
