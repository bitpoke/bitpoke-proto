import * as React from 'react'
import { Field } from 'redux-form'

import { get } from 'lodash'

import { forms, api, sites } from '../redux'

import { withForm } from '../components/Form'
import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import SiteTitle from '../components/SiteTitle'

type Props = forms.Props<sites.ISitePayload>

const SiteForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const entry = get(initialValues, 'site')

    return (
        <FormContainer
            title={ <SiteTitle entry={ entry } /> }
            { ...props }
        >
            <Field
                name="site.primaryDomain"
                label="Domain Name"
                component={ InputField }
            />
        </FormContainer>
    )
}

export default withForm({ name: forms.Name.site })(SiteForm)
