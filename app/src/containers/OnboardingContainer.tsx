import React, { Fragment } from 'react'

import { Switch, Route } from 'react-router-dom'

import Container from '../components/Container'
import OrganizationForm from '../components/OrganizationForm'
import ProjectForm from '../components/ProjectForm'

import { routing } from '../redux'

const OnboardingContainer: React.SFC<{}> = () => {
    return (
        <Container>
            <Switch>
                <Route
                    path={ routing.routeFor('onboarding', { step: 'project' }) }
                    component={ ProjectForm }
                />
                <Route
                    path={ routing.routeFor('onboarding') }
                    component={ OrganizationForm }
                />
            </Switch>
        </Container>
    )
}

export default OnboardingContainer
