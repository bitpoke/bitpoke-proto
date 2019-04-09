import * as React from 'react'
import { connect } from 'react-redux'
import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState, DispatchProp, routing, sites } from '../redux'

import TitleBar from '../components/TitleBar'
import SitesList from '../components/SitesList'

type OwnProps = {
    entry: sites.ISite | null
}

type ReduxProps = {
    isEditing: boolean
}

type Props = OwnProps & ReduxProps & DispatchProp

const SiteTitle: React.SFC<Props> = ({ entry, isEditing, dispatch }) => {
    if (!entry) {
        return null
    }

    return (
        <TitleBar
            title={ entry.primaryDomain }
            subtitle={ entry.name }
            actions={
                isEditing
                    ? (
                        <ButtonGroup>
                            <Button
                                text="Discard"
                                icon="cross"
                                intent={ Intent.PRIMARY }
                                onClick={ () =>
                                    dispatch(routing.push(
                                        routing.routeForResource(entry)
                                    ))
                                }
                            />
                        </ButtonGroup>
                    ) : (
                        <ButtonGroup>
                            <Button
                                text="Edit site"
                                icon="edit"
                                intent={ Intent.PRIMARY }
                                onClick={ () =>
                                    dispatch(routing.push(
                                        routing.routeForResource(entry, { action: 'edit' })
                                    ))
                                }
                            />
                            <Button
                                text="Delete site"
                                icon="trash"
                                intent={ Intent.DANGER }
                                onClick={ () => dispatch(sites.destroy(entry)) }
                            />
                        </ButtonGroup>
                    )
            }
        />
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    const currentRoute = routing.getCurrentRoute(state)
    const isEditing = currentRoute.key === 'sites' && get(currentRoute, 'params.action') === 'edit'
    return {
        isEditing
    }
}

export default connect(mapStateToProps)(SiteTitle)
