import * as React from 'react'
import { connect } from 'react-redux'
import { Switch, Route } from 'react-router-dom'

import Container from '../components/Container'
import OrganizationForm from '../components/OrganizationForm'
import OrganizationDetails from '../components/OrganizationDetails'

import { RootState, routing, organizations } from '../redux'

type ReduxProps = {
    organization: organizations.IOrganization | null
}

type Props = ReduxProps

const OrganizationsContainer: React.SFC<Props> = ({ organization }) => {
    return (
        <Container>
            <Switch>
                <Route
                    path={ routing.routeFor('organization', { action: 'new' }) }
                    component={ OrganizationForm }
                />
                { organization && (
                    <Route
                        path={ routing.routeForResource(organization, { action: 'edit' }) }
                        render={ () => <OrganizationForm initialValues={ { organization } } /> }
                    />
                ) }
                <Route
                    path={ routing.routeFor('organization') }
                    render={ () => <OrganizationDetails entry={ organization } /> }
                />
            </Switch>
        </Container>
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    return {
        organization: organizations.getForCurrentURL(state)
    }
}

export default connect(mapStateToProps)(OrganizationsContainer)
