import * as React from 'react'
import { Field } from 'redux-form'
import { get } from 'lodash'

import { forms, api, projects } from '../redux'

import { withForm } from '../components/Form'
import FormContainer from '../components/FormContainer'
import InputField from '../components/InputField'
import TitleBar from '../components/TitleBar'
import ProjectTitle from '../components/ProjectTitle'

type Props = forms.Props<projects.IProjectPayload>

const ProjectForm: React.SFC<Props> = (props) => {
    const { initialValues } = props

    const entry = get(initialValues, 'project', {})
    const isNewEntry = api.isNewEntry(entry)
    const title = isNewEntry
        ? <TitleBar title="Create Project" />
        : <ProjectTitle entry={ entry } />

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
