import * as React from 'react'
import { connect } from 'react-redux'
import { singular } from 'pluralize'
import { Button, ButtonGroup, Intent } from '@blueprintjs/core'

import { get, isFunction } from 'lodash'

import { RootState, DispatchProp, api, routing } from '../redux'

import DestroyButton from '../components/DestroyButton'


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
        entry, dispatch,
        onCreate, onGenerate, onDestroy,
        minimal, withTitles, isEditing, isCreating, isHidden
    } = props

    const resourceName = singular(props.resourceName)

    if (isHidden) {
        return null
    }

    if (isCreating) {
        return (
            <Button
                text={ withTitles && `Discard ${resourceName}` }
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
                        text={ withTitles && `Create ${resourceName}` }
                        icon="add"
                        intent={ Intent.SUCCESS }
                        minimal={ minimal }
                        onClick={ () => onCreate() }
                    />
                ) }
                { isFunction(onGenerate) && (
                    <Button
                        text={ withTitles && `Generate random ${resourceName}` }
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
                text={ withTitles && `Edit ${resourceName}` }
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
                <DestroyButton
                    text={ withTitles && `Delete ${resourceName}` }
                    minimal={ minimal }
                    confirmationText={
                        `Are you sure you want to delete this ${resourceName}?\nThis action is not undoable!`
                    }
                    onDestroy={ () => onDestroy() }
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
