import * as React from 'react'

import { Field } from 'redux-form'

import { forms, api, sites } from '../redux'

import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import { withForm } from '../components/Form'

type Props = forms.Props<sites.ISitePayload>

const SiteForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const title = api.isNewEntry(initialValues)
        ? 'Create Site'
        : initialValues.primaryDomain

    return (
        <FormContainer title={ title } { ...props }>
            <Field
                name="site.primaryDomain"
                label="Domain Name"
                component={ InputField }
            />
        </FormContainer>
    )
}

export default withForm({ name: forms.Name.site })(SiteForm)
