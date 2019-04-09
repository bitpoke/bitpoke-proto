import * as React from 'react'
import { Field } from 'redux-form'
import { get } from 'lodash'

import { forms, api, projects } from '../redux'

import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import { withForm } from '../components/Form'

type Props = forms.Props<projects.IProjectPayload>

const ProjectForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const entry = get(initialValues, 'project', {})
    const title = api.isNewEntry(entry)
        ? 'Create Project'
        : entry.displayName

    return (
        <FormContainer title={ title } { ...props }>
            <Field
                name="project.displayName"
                label="Name"
                component={ InputField }
            />
        </FormContainer>
    )
}

export default withForm({ name: forms.Name.project })(ProjectForm)
