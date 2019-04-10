import * as React from 'react'
import { connect } from 'react-redux'
import { singular } from 'pluralize'
import { Button, ButtonGroup, Card, Elevation, Intent } from '@blueprintjs/core'

import { get, isFunction } from 'lodash'

import { RootState, DispatchProp, api, routing } from '../redux'

import TitleBar from '../components/TitleBar'
import SitesList from '../components/SitesList'

type OwnProps = {
    resourceName: api.Resource,
    entry?: api.AnyResourceInstance | null,
    withTitles?: boolean,
    minimal?: boolean,
    onCreate?: () => void,
    onGenerate?: () => void,
    onDestroy?: () => void
}

type ReduxProps = {
    isEditing: boolean,
    isCreating: boolean,
    isHidden: boolean
}

type Props = OwnProps & ReduxProps & DispatchProp

const ResourceActions: React.SFC<Props> = (props) => {
    const {
        entry, resourceName, dispatch,
        onCreate, onGenerate, onDestroy,
        minimal, withTitles, isEditing, isCreating, isHidden
    } = props

    if (isHidden) {
        return null
    }

    if (isCreating) {
        return (
            <Button
                text={ withTitles && `Discard ${singular(resourceName)}` }
                icon="cross"
                intent={ Intent.PRIMARY }
                minimal={ minimal }
                onClick={ () => dispatch(routing.goBack()) }
            />
        )
    }

    if (!entry && !isCreating && !isEditing) {
        return (
            <ButtonGroup>
                { isFunction(onCreate) && (
                    <Button
                        text={ withTitles && `Create ${singular(resourceName)}` }
                        icon="add"
                        intent={ Intent.SUCCESS }
                        minimal={ minimal }
                        onClick={ () => onCreate() }
                    />
                ) }
                { isFunction(onGenerate) && (
                    <Button
                        text={ withTitles && `Generate random ${singular(resourceName)}` }
                        icon="random"
                        intent={ Intent.SUCCESS }
                        minimal={ minimal }
                        onClick={ () => onGenerate() }
                    />
                ) }
            </ButtonGroup>
        )
    }

    if (!entry) {
        return null
    }

    if (isEditing) {
        return (
            <Button
                text={ withTitles && 'Discard changes' }
                icon="cross"
                intent={ Intent.PRIMARY }
                minimal={ minimal }
                onClick={ () =>
                    dispatch(routing.goBack())
                }
            />
        )
    }

    return (
        <ButtonGroup>
            <Button
                text={ withTitles && `Edit ${singular(resourceName)}` }
                icon="edit"
                intent={ Intent.PRIMARY }
                minimal={ minimal }
                onClick={ (e: React.SyntheticEvent<EventTarget>) => {
                    e.stopPropagation()
                    dispatch(routing.push(
                        routing.routeForResource(entry, { action: 'edit' })
                    ))
                } }
            />
            { isFunction(onDestroy) && (
                <Button
                    text={ withTitles && `Delete ${singular(resourceName)}` }
                    icon="trash"
                    intent={ Intent.DANGER }
                    minimal={ minimal }
                    onClick={ (e: React.SyntheticEvent<EventTarget>) => {
                        e.stopPropagation()
                        onDestroy()
                    } }
                />
            ) }
        </ButtonGroup>
    )
}

ResourceActions.defaultProps = {
    withTitles: true,
    minimal: false
}

function mapStateToProps(state: RootState, ownProps: OwnProps): ReduxProps {
    const { resourceName } = ownProps
    const currentRoute = routing.getCurrentRoute(state)
    const action = get(currentRoute, 'params.action')
    const isEditing =  action === 'edit'
    const isCreating = action === 'new'
    const isHidden = (isEditing || isCreating) && currentRoute.key !== resourceName
    return {
        isEditing,
        isCreating,
        isHidden
    }
}

export default connect(mapStateToProps)(ResourceActions)
