import * as React from 'react'
import {  } from 'redux'
import { connect } from 'react-redux'
import {
    reduxForm,
    getFormValues,
    getFormSubmitErrors,
    ConfigProps,
    InjectedFormProps
} from 'redux-form'

import { omit, get } from 'lodash'

import { RootState, DispatchProp, forms } from '../redux'

import { Omit } from '../utils'

type Config = Omit<ConfigProps<any, any, string>, 'form'> & { name: forms.Name }

export function withForm(config: Config) {
    const { name } = config

    return (WrappedComponent: React.ComponentType<any>) => {
        const Form = (props: InjectedFormProps & DispatchProp) => {
            const {
                initialValues, handleSubmit,
                submitting, dirty, pristine, valid, error,
                dispatch, ...otherProps
            } = props

            const onSubmit = handleSubmit((values) => (
                new Promise((resolve, reject) => {
                    dispatch(forms.submit({ name, values, resolve, reject }))
                })
            ))

            const newProps = omit({
                ...otherProps,
                isSubmitting  : get(props, 'submitting', false),
                isDirty       : get(props, 'dirty', false),
                isPristine    : get(props, 'pristine', false),
                isValid       : get(props, 'valid', false),
                initialValues : get(props, 'initialValues', {}),
                onSubmit,
                error
            }, ['submitting', 'dirty', 'pristine', 'valid'])

            return (
                <form onSubmit={ onSubmit }>
                    <WrappedComponent { ...newProps } />
                </form>
            )
        }

        function mapStateToProps(state: RootState) {
            const currentValues = getFormValues(name, forms.getState)(state) || {}
            const errors = getFormSubmitErrors(name, forms.getState)(state) || {}
            return {
                currentValues,
                errors
            }
        }

        const ReduxForm = reduxForm({
            form               : name,
            enableReinitialize : true,
            destroyOnUnmount   : false,
            getFormState       : forms.getState,
            ...omit(config, 'name')
        })(Form)

        return connect(mapStateToProps)(ReduxForm)
    }
}
