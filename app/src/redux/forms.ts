import {
    reducer, reset, SubmissionError, isSubmitting,
    FormStateMap, InjectedFormProps, WrappedFieldProps
} from 'redux-form'
import { ActionType, createAsyncAction, action as createAction, isOfType } from 'typesafe-actions'
import { Dispatch } from 'redux'
import { takeEvery, put } from 'redux-saga/effects'

import { toUpper, snakeCase } from 'lodash'

import { RootState } from '../redux'
import { Omit } from '../utils'

export enum Name {
    organization = 'organization',
    project = 'project'
}

export type Values = {}

export type Payload = {
    name: Name,
    values: Values,
    resolve: () => void
    reject: (error?: any) => void
}

export type Actions = ActionType<typeof actions>
export type State = FormStateMap

export { SubmissionError, reset }

type DefaultInjectedProps<FormData = {}> = InjectedFormProps<FormData, {}, string>
type UpdatedInjectedProps<FormData = {}> =
    Omit<DefaultInjectedProps<FormData>, 'submitted' | 'dirty' | 'pristine' | 'valid'> & {
        isSubmitting: boolean,
        isDirty: boolean,
        isPristine: boolean,
        isValid: boolean
    }

export type Props<FormData = {}> = UpdatedInjectedProps<FormData> & {
    dispatch: Dispatch
    onSubmit: () => void
}

export type FieldProps = WrappedFieldProps & {
    label?: string
    helperText?: string
}

//
//  ACTIONS

export const SUBMITTED = '@ forms / SUBMITTED'

export const submit = (payload: Payload) => createAction(SUBMITTED, payload)

const actions = {
    submit
}


//
//  REDUCER

export { reducer }


//
//  SAGA

export function* saga() {
    yield takeEvery(SUBMITTED, emitFormActions)
}

function* emitFormActions(action: ActionType<typeof submit>) {
    const { payload } = action
    const { name } = payload

    const type = createDescriptor(name)

    yield put({ type, payload })
}

export function* takeEverySubmission(name: Name, handler: (action: any) => IterableIterator<any>) {
    const descriptor = createDescriptor(name)
    yield takeEvery(descriptor, handler)
}

//
//  HELPERS and UTILITIES

function createDescriptor(name: Name) {
    return `@ forms / ${toUpper(snakeCase(name))}_FORM_SUBMITTED`
}

//
//  SELECTORS

export const getState = (state: RootState) => state.forms
