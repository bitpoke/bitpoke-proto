import React from 'react'
import Container from '../components/Container'
import ProjectsList from '../components/ProjectsList'

const DashboardContainer: React.SFC<{}> = () => {
    return (
        <Container>
            <ProjectsList />
        </Container>
    )
}

export default DashboardContainer
