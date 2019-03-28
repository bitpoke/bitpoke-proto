import * as React from 'react'

import { has } from 'lodash'
import { Field } from 'redux-form'

import { forms, api, organizations } from '../redux'

import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import { withForm } from '../components/Form'

type Props = forms.Props<organizations.IOrganization>

const OrganizationForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const title = api.isNewEntry(initialValues)
        ? 'Create Organization'
        : initialValues.displayName

    return (
        <FormContainer title={ title } { ...props }>
            <Field
                name="displayName"
                label="Name"
                component={ InputField }
            />
        </FormContainer>
    )
}

export default withForm({ name: forms.Name.organization })(OrganizationForm)
