import * as React from 'react'
import { Field } from 'redux-form'
import { get } from 'lodash'

import { forms, api, organizations } from '../redux'

import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import { withForm } from '../components/Form'

type Props = forms.Props<organizations.IOrganizationPayload>

const OrganizationForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const entry = get(initialValues, 'organization', {})
    const title = api.isNewEntry(entry)
        ? 'Create Organization'
        : entry.displayName

    return (
        <FormContainer title={ title } { ...props }>
            <Field
                name="organization.displayName"
                label="Name"
                component={ InputField }
            />
        </FormContainer>
    )
}

export default withForm({ name: forms.Name.organization })(OrganizationForm)
