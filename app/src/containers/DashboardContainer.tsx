import React, { Fragment } from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import { map, get } from 'lodash'

import { Button, Intent } from '@blueprintjs/core'

import { RootState, api, organizations } from '../redux'

import Container from '../components/Container'
import ProjectsList from '../components/ProjectsList'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    currentOrganization: organizations.IOrganization | null
}

const DashboardContainer: React.SFC<Props & ReduxProps> = ({ currentOrganization, dispatch }) => {
    if (!currentOrganization) {
        return null
    }

    return (
        <Container>
            <h2>{ currentOrganization.displayName }</h2>
            <p>{ currentOrganization.name }</p>
            <Button
                text="Delete Organization"
                icon="trash"
                intent={ Intent.DANGER }
                onClick={ () => dispatch(organizations.destroy(currentOrganization)) }
            />
            <ProjectsList />
        </Container>
    )
}

function mapStateToProps(state: RootState) {
    const currentOrganization = organizations.getCurrent(state)
    return {
        currentOrganization
    }
}

export default connect(mapStateToProps)(DashboardContainer)
