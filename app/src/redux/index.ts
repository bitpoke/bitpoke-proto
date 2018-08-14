import * as routing from './routing'
import * as projects from './projects'

export type RootState = {
    routing  : routing.State,
    projects : projects.State
}
export type AnyAction =
    auth.Actions
    | routing.Actions
    | projects.Actions

export {
    routing,
    projects
}
