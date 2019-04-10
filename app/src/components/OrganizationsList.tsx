import React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { map, get } from 'lodash'

import { Navbar, Alignment, Button, Intent } from '@blueprintjs/core'

import { RootState, DispatchProp, organizations, api, routing } from '../redux'

type ReduxProps = {
    entries: api.ResourcesList<organizations.IOrganization>,
    selectedEntry: organizations.IOrganization | null
}

type Props = ReduxProps & DispatchProp

const { Group, Divider } = Navbar

const OrganizationsList: React.SFC<Props> = ({ entries, selectedEntry, dispatch }) => {
    return (
        <Group align={ Alignment.LEFT }>
            <Divider />
            { map(entries, (organization) => (
                <Button
                    minimal
                    active={ organization === selectedEntry }
                    key={ organization.name }
                    text={ organization.displayName }
                    onClick={ () => {
                        dispatch(routing.push(
                            routing.routeFor('dashboard', { org: organizations.parseName(organization.name).slug })
                        ))
                    } }
                />
            )) }
            <Button
                minimal
                intent={ Intent.SUCCESS }
                icon="add"
                onClick={ () => {
                    dispatch(organizations.create({ organization: { displayName: faker.company.companyName() } }))
                } }
            />
        </Group>
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    const selectedEntry = organizations.getCurrent(state)
    const entries = organizations.getAll(state)

    return {
        entries,
        selectedEntry
    }
}

export default connect(mapStateToProps)(OrganizationsList)
