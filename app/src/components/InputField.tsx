import * as React from 'react'

import { InputGroup, Intent } from '@blueprintjs/core'

import { forms } from '../redux'

import FormGroup from '../components/FormGroup'

type Props = forms.FieldProps

const InputField: React.SFC<Props> = (props) => {
    const { input, meta, label, ...otherProps } = props
    const { name } = input
    const { error } = meta

    return (
        <FormGroup { ...props }>
            <InputGroup
                id={ name }
                { ...otherProps }
                { ...input }
                intent={ error ? Intent.DANGER : Intent.PRIMARY }
            />
        </FormGroup>
    )
}

export default InputField
