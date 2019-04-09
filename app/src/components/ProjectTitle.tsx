import * as React from 'react'
import { connect } from 'react-redux'
import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { get } from 'lodash'

import { RootState, DispatchProp, routing, projects } from '../redux'

import TitleBar from '../components/TitleBar'
import SitesList from '../components/SitesList'

type OwnProps = {
    entry: projects.IProject | null
}

type ReduxProps = {
    isEditing: boolean
}

type Props = OwnProps & ReduxProps & DispatchProp

const ProjectTitle: React.SFC<Props> = ({ entry, isEditing, dispatch }) => {
    if (!entry) {
        return null
    }

    return (
        <TitleBar
            title={ entry.displayName }
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
                                text="Edit project"
                                icon="edit"
                                intent={ Intent.PRIMARY }
                                onClick={ () =>
                                    dispatch(routing.push(
                                        routing.routeForResource(entry, { action: 'edit' })
                                    ))
                                }
                            />
                            <Button
                                text="Delete project"
                                icon="trash"
                                intent={ Intent.DANGER }
                                onClick={ () => dispatch(projects.destroy(entry)) }
                            />
                        </ButtonGroup>
                    )
            }
        />
    )
}

function mapStateToProps(state: RootState): ReduxProps {
    const currentRoute = routing.getCurrentRoute(state)
    const isEditing = currentRoute.key === 'projects' && get(currentRoute, 'params.action') === 'edit'
    return {
        isEditing
    }
}

export default connect(mapStateToProps)(ProjectTitle)

