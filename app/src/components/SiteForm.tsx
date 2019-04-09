import * as React from 'react'
import { Field } from 'redux-form'

import { get } from 'lodash'

import { forms, api, sites } from '../redux'

import { withForm } from '../components/Form'
import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import TitleBar from '../components/TitleBar'
import SiteTitle from '../components/SiteTitle'

type Props = forms.Props<sites.ISitePayload>

const SiteForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const entry = get(initialValues, 'site', {})
    const isNewEntry = api.isNewEntry(entry)
    const title = isNewEntry
        ? <TitleBar title="Create Site" />
        : <SiteTitle entry={ entry } />

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
