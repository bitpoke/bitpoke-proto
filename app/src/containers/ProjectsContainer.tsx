import * as React from 'react'
import { connect } from 'react-redux'
import { Switch, Route } from 'react-router-dom'

import { isEqual } from 'lodash'

import Container from '../components/Container'
import ProjectTitle from '../components/ProjectTitle'
import ProjectForm from '../components/ProjectForm'
import ProjectDetails from '../components/ProjectDetails'
import SiteForm from '../components/SiteForm'
import SiteDetails from '../components/SiteDetails'

import { RootState, DispatchProp, routing, organizations, projects } from '../redux'

type ReduxProps = {
    project: projects.IProject | null,
    currentOrganization: organizations.IOrganization | null,
    currentRoute: routing.Route
}

type Props = ReduxProps & DispatchProp

const ProjectsContainer: React.SFC<Props> = (props) => {
    const { project } = props

    if (!project) {
        return (
            <Container>
                <ProjectForm initialValues={ {} } />
            </Container>
        )
    }

    return (
        <Container>
            <Switch>
                <Route
                    path={ routing.routeForResource(project, { action: 'new-site' }) }
                    render={ () => (
                        <div>
                            <ProjectTitle entry={ project } />
                            <SiteForm initialValues={ { parent: project.name } } />
                        </div>
                    ) }
                />
                <Route
                    path={ routing.routeForResource(project, { action: 'edit' }) }
                    render={ () => <ProjectForm initialValues={ { project } } /> }
                />
                <Route
                    path={ routing.routeForResource(project, { action: 'edit' }) }
                    render={ () => <ProjectForm initialValues={ { project } } /> }
                />
                <Route
                    path={ routing.routeForResource(project) }
                    render={ () => <ProjectDetails entry={ project } /> }
                />
            </Switch>
        </Container>
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    return {
        currentRoute: routing.getCurrentRoute(state),
        currentOrganization: organizations.getCurrent(state),
        project: projects.getForCurrentURL(state)
    }
}

export default connect(mapStateToProps)(ProjectsContainer)
