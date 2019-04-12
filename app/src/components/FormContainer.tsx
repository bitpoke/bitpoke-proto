import * as React from 'react'

import { isFunction } from 'lodash'

import { Card, Elevation, Button, Intent } from '@blueprintjs/core'

import { forms } from '../redux'

type Props = forms.Props & {
    title?: React.ReactNode
}

const FormContainer: React.SFC<Props> = (props) => {
    const { title, onSubmit, isSubmitting, children } = props

    return (
        <div>
            { title }
            <Card elevation={ Elevation.TWO }>
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
        </div>
    )
}

export default FormContainer
