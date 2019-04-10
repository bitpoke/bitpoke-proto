import * as React from 'react'
import { Field } from 'redux-form'

import { get } from 'lodash'

import { forms, organizations } from '../redux'

import { withForm } from '../components/Form'
import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import OrganizationTitle from '../components/OrganizationTitle'

type Props = forms.Props<organizations.IOrganizationPayload>

const OrganizationForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const entry = get(initialValues, 'organization')

    return (
        <FormContainer
            title={ <OrganizationTitle entry={ entry } /> }
            { ...props }
        >
            <Field
                name="organization.displayName"
                label="Name"
                component={ InputField }
            />
        </FormContainer>
    )
}

export default withForm({ name: forms.Name.organization })(OrganizationForm)
