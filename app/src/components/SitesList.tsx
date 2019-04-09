import * as React from 'react'
import { connect } from 'react-redux'
import faker from 'faker'

import { map } from 'lodash'

import { Button, Card, Elevation, Intent } from '@blueprintjs/core'

import { RootState, DispatchProp, api, routing, sites, projects } from '../redux'

type OwnProps = {
    project : projects.ProjectName
}

type ReduxProps = {
    entries: api.ResourcesList<sites.ISite>
}

type Props = OwnProps & ReduxProps & DispatchProp

const SitesList: React.SFC<Props> = (props) => {
    const { entries, project, dispatch } = props
    return (
        <div>
            <h2>Sites</h2>
            <Button
                text="Create site"
                icon="add"
                intent={ Intent.SUCCESS }
                onClick={ () => dispatch(routing.push(routing.routeFor('site', { step: 'new' }))) }
            />
            <Button
                text="Create random site"
                icon="random"
                intent={ Intent.SUCCESS }
                onClick={ () => dispatch(sites.create({
                    parent: project,
                    site: {
                        primaryDomain: faker.internet.domainName()
                    }
                })) }
            />
            { map(entries, (site) => (
                <Card
                    key={ site.name }
                    elevation={ Elevation.TWO }
                >
                    <h5><a href="#">{ site.primaryDomain }</a></h5>
                    <p>{ site.primaryDomain }</p>
                    <Button
                        text="Delete site"
                        icon="trash"
                        intent={ Intent.DANGER }
                        onClick={ () => dispatch(sites.destroy(site)) }
                    />
                </Card>
            )) }
        </div>
    )
}

function mapStateToProps(state: RootState, ownProps: OwnProps): ReduxProps {
    const entries = sites.getForProject(ownProps.project)(state)
    return {
        entries
    }
}

export default connect(mapStateToProps)(SitesList)
