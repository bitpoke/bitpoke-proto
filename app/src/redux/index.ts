import * as app from './app'
import * as auth from './auth'
import * as routing from './routing'
import * as projects from './projects'

export type RootState = {
    app      : app.State,
    auth     : auth.State,
    routing  : routing.State,
    projects : projects.State
}
export type AnyAction =
    app.Actions
    | auth.Actions
    | routing.Actions
    | projects.Actions

export {
    app,
    auth,
    routing,
    projects
}
