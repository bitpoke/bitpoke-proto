import * as React from 'react'
import { connect } from 'react-redux'

import { RootState, organizations } from '../redux'

import Container from '../components/Container'
import ProjectsList from '../components/ProjectsList'

type ReduxProps = {
    currentOrganization: organizations.IOrganization | null
}

const DashboardContainer: React.SFC<ReduxProps> = ({ currentOrganization }) => {
    if (!currentOrganization) {
        return null
    }

    return (
        <Container>
            <ProjectsList organization={ currentOrganization.name } />
        </Container>
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    const currentOrganization = organizations.getCurrent(state)

    return {
        currentOrganization
    }
}

export default connect(mapStateToProps)(DashboardContainer)
