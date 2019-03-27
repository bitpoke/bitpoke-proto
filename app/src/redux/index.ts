import * as app from './app'
import * as api from './api'
import * as grpc from './grpc'
import * as auth from './auth'
import * as routing from './routing'
import * as organizations from './organizations'
import * as projects from './projects'

export type RootState = {
    app           : app.State,
    grpc          : grpc.State,
    auth          : auth.State,
    routing       : routing.State,
    organizations : organizations.State,
    projects      : projects.State
}

export type Reducer = (state: RootState | undefined, action: AnyAction) => RootState

export type AnyAction =
    app.Actions
    | grpc.Actions
    | auth.Actions
    | routing.Actions
    | organizations.Actions
    | projects.Actions

export {
    app,
    api,
    grpc,
    auth,
    routing,
    organizations,
    projects
}
