import React, { Fragment } from 'react'
import { Dispatch } from 'redux'
import { connect } from 'react-redux'

import faker from 'faker'

import { map } from 'lodash'

import { Navbar, Alignment, Button, Intent } from '@blueprintjs/core'

import { RootState, organizations, api } from '../redux'

type Props = {
    dispatch: Dispatch
}

type ReduxProps = {
    entries: api.ResourcesList<organizations.IOrganization>,
    selectedEntry: organizations.IOrganization | null
}

const { Group, Heading, Divider } = Navbar

const OrganizationsList: React.SFC<Props & ReduxProps> = ({ entries, selectedEntry, dispatch }) => {
    return (
        <Group align={Alignment.LEFT}>
            <Divider />
            { map(entries, (organization) => (
                <Button
                    minimal
                    active={ organization === selectedEntry }
                    key={ `organization-${organization.name}` }
                    text={ organization.displayName }
                    onClick={ () => {
                        dispatch(organizations.select(organization))
                    } }
                />
            )) }
            <Button
                minimal
                intent={ Intent.SUCCESS }
                text="Create random organization"
                onClick={ () => {
                    dispatch(organizations.create({ displayName: faker.company.companyName() }))
                } }
            />
        </Group>
    )
}

function mapStateToProps(state: RootState) {
    const selectedEntry = organizations.getCurrent(state)
    const entries = organizations.getAll(state)

    return {
        entries,
        selectedEntry
    }
}

export default connect(mapStateToProps)(OrganizationsList)
