import * as React from 'react'

import { isFunction } from 'lodash'

import { Card, Elevation, Button, Intent } from '@blueprintjs/core'

import { forms } from '../redux'

type Props = forms.Props & {
    title?: string | null
}

const FormContainer: React.SFC<Props> = (props) => {
    const { title, onSubmit, isSubmitting, children } = props

    return (
        <Card elevation={ Elevation.TWO }>
            { title && <h2>{ title }</h2> }
            { children }
            { isFunction(onSubmit) && (
                <Button
                    text="Save"
                    onClick={ onSubmit }
                    intent={ Intent.SUCCESS }
                    loading={ isSubmitting }
                />
            ) }
        </Card>
    )
}

export default FormContainer
