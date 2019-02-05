import * as app from './app'
import * as auth from './auth'
import * as routing from './routing'
import * as organizations from './organizations'

export type RootState = {
    app           : app.State,
    auth          : auth.State,
    routing       : routing.State,
    organizations : organizations.State
}
export type AnyAction =
    app.Actions
    | auth.Actions
    | routing.Actions
    | organizations.Actions

export {
    app,
    auth,
    routing,
    organizations
}
