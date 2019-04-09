import * as React from 'react'

import { FormGroup as BaseFormGroup, Intent } from '@blueprintjs/core'

import { forms } from '../redux'

type Props = forms.FieldProps

const FormGroup: React.SFC<Props> = (props) => {
    const { input, meta } = props
    const { name } = input
    const { error } = meta

    return (
        <BaseFormGroup
            labelFor={ name }
            { ...props }
            helperText={ error ? error : props.helperText }
            intent={ error ? Intent.DANGER : Intent.PRIMARY }
        />
    )
}

export default FormGroup
