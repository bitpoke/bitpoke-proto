import { Dispatch } from 'redux'

import * as app from './app'
import * as api from './api'
import * as grpc from './grpc'
import * as auth from './auth'
import * as forms from './forms'
import * as routing from './routing'
import * as organizations from './organizations'
import * as projects from './projects'
import * as sites from './sites'
import toasts from './toasts'

export type RootState = {
    app           : app.State,
    grpc          : grpc.State,
    auth          : auth.State,
    forms         : forms.State
    routing       : routing.State,
    organizations : organizations.State,
    projects      : projects.State,
    sites         : sites.State
}

export type Reducer = (state: RootState | undefined, action: AnyAction) => RootState

export type AnyAction =
    app.Actions
    | grpc.Actions
    | auth.Actions
    | forms.Actions
    | routing.Actions
    | organizations.Actions
    | projects.Actions
    | sites.Actions

export type ActionDescriptor = string

export type DispatchProp = { dispatch: Dispatch }

export {
    app,
    api,
    grpc,
    auth,
    forms,
    routing,
    organizations,
    projects,
    sites,
    toasts
}
