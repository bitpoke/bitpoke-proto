import React, { Fragment } from 'react'

import Container from '../components/Container'
import OrganizationForm from '../components/OrganizationForm'

const OnboardingContainer: React.SFC<{}> = () => {
    return (
        <Container>
            <OrganizationForm />
        </Container>
    )
}

export default OnboardingContainer
