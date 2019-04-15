import { expectSaga } from 'redux-saga-test-plan'
import { map, reduce, keys, values, size, isArray } from 'lodash'

import { api, grpc, sites } from '../'

describe('api', () => {
    const { createActionDescriptor } = api
    describe('createActionDescriptor()', () => {
        it('builds action descriptors for resources based on request type and status', () => {
            expect(createActionDescriptor(api.Request.list, api.Status.succeeded)).toEqual('LIST_SUCCEEDED')
            expect(createActionDescriptor(api.Request.create, api.Status.requested)).toEqual('CREATE_REQUESTED')
            expect(createActionDescriptor(api.Request.update, api.Status.failed)).toEqual('UPDATE_FAILED')
        })
    })

    const { createActionTypes } = api
    describe('createActionTypes()', () => {
        it('builds action types for all request types and status combinations', () => {
            expect(createActionTypes(api.Resource.site)).toEqual({
                LIST_REQUESTED    : '@ sites / LIST_REQUESTED',
                LIST_SUCCEEDED    : '@ sites / LIST_SUCCEEDED',
                LIST_FAILED       : '@ sites / LIST_FAILED',

                GET_REQUESTED     : '@ sites / GET_REQUESTED',
                GET_SUCCEEDED     : '@ sites / GET_SUCCEEDED',
                GET_FAILED        : '@ sites / GET_FAILED',

                CREATE_REQUESTED  : '@ sites / CREATE_REQUESTED',
                CREATE_SUCCEEDED  : '@ sites / CREATE_SUCCEEDED',
                CREATE_FAILED     : '@ sites / CREATE_FAILED',

                UPDATE_REQUESTED  : '@ sites / UPDATE_REQUESTED',
                UPDATE_SUCCEEDED  : '@ sites / UPDATE_SUCCEEDED',
                UPDATE_FAILED     : '@ sites / UPDATE_FAILED',

                DESTROY_REQUESTED : '@ sites / DESTROY_REQUESTED',
                DESTROY_SUCCEEDED : '@ sites / DESTROY_SUCCEEDED',
                DESTROY_FAILED    : '@ sites / DESTROY_FAILED'
            })
        })
    })

    const { createReducer } = api
    describe('reducer()', () => {
        const actionTypes = createActionTypes(api.Resource.site)
        const reducer = createReducer(api.Resource.site, actionTypes)

        const createEntry = (name, otherProps = {}) => ({ name, ...otherProps })
        const createState = (entries) => ({ entries })
        const createResponse = (payload) => {
            const data = isArray(payload)
                ? { [api.Resource.site]: payload }
                : payload

            return {
                data,
                error: null,
                request: {}
            }
        }

        const createDestroyResponse = (payload) => ({
            data: null,
            error: null,
            request: {
                data: payload
            }
        })

        const createResponseAction = (type, entry) => {
            const payload = type === actionTypes.DESTROY_SUCCEEDED
                ? createDestroyResponse(entry)
                : createResponse(entry)

            return {
                type,
                payload
            }
        }

        const existingEntry = createEntry('proj/abc/sites/x', { primaryDomain: 'x.com' })

        const initialState = createState({})
        const existingState = createState({ [existingEntry.name]: existingEntry })

        describe('handles LIST requests', () => {
            const fetchedEntries = [
                createEntry('proj/abc/sites/a'),
                createEntry('proj/abc/sites/b'),
                createEntry('proj/abc/sites/c')
            ]

            it('storing fetched entries indexed by name', () => {
                const action = createResponseAction(actionTypes.LIST_SUCCEEDED, fetchedEntries)
                const state = reducer(initialState, action)
                expect(keys(state.entries)).toEqual(map(fetchedEntries, 'name'))
                expect(values(state.entries)).toEqual(fetchedEntries)
            })

            it('merging existing entries with fetched ones', () => {
                const action = createResponseAction(actionTypes.LIST_SUCCEEDED, fetchedEntries)
                const state = reducer(existingState, action)

                expect(keys(state.entries)).toEqual([existingEntry.name, ...map(fetchedEntries, 'name')])
                expect(values(state.entries)).toEqual([existingEntry, ...fetchedEntries])
            })

            it('ignoring invalid actions/payloads', () => {
                expect(
                    reducer(initialState, createResponseAction(actionTypes.LIST_SUCCEEDED, null))
                ).toEqual(initialState)

                expect(
                    reducer(initialState, createResponseAction(actionTypes.LIST_FAILED, fetchedEntries))
                ).toEqual(initialState)

                expect(
                    reducer(initialState, createResponseAction(actionTypes.GET_SUCCEEDED, fetchedEntries))
                ).toEqual(initialState)
            })
        })

        describe('handles GET, CREATE, UPDATE requests', () => {
            it('merging existing entries with fetched entry payload', () => {
                const fetchedEntry = createEntry('proj/abc/sites/y')

                const action = createResponseAction(actionTypes.GET_SUCCEEDED, fetchedEntry)
                const state = reducer(existingState, action)

                expect(keys(state.entries)).toEqual([existingEntry.name, fetchedEntry.name])
                expect(values(state.entries)).toEqual([existingEntry, fetchedEntry])
            })

            it('replacing keys for existing payload with updated entry payload', () => {
                const updatedEntry = createEntry('proj/abc/sites/x', { primaryDomain: 'y.com' })

                const action = createResponseAction(actionTypes.UPDATE_SUCCEEDED, updatedEntry)
                const state = reducer(existingState, action)

                expect(state.entries[existingEntry.name].primaryDomain).toEqual('y.com')
            })

            it('ignoring invalid actions/payloads', () => {
                expect(
                    reducer(initialState, createResponseAction(actionTypes.CREATE_SUCCEEDED, {}))
                ).toEqual(initialState)

                expect(
                    reducer(initialState, createResponseAction(actionTypes.CREATE_FAILED, {}))
                ).toEqual(initialState)
            })
        })

        describe('handles DESTROY requests', () => {
            it('removing the destroyed entry from the state', () => {
                const action = createResponseAction(actionTypes.DESTROY_SUCCEEDED, existingEntry)
                const state = reducer(existingState, action)

                expect(state).toEqual(initialState)
            })

            it('ignoring invalid actions/payloads', () => {
                expect(
                    reducer(existingState, createResponseAction(actionTypes.DESTROY_FAILED, existingEntry))
                ).toEqual(existingState)
            })
        })
    })

    const { createSelectors } = api
    describe('selectors', () => {
        const { getState, getAll, countAll, getByName, getForURL } = createSelectors(api.Resource.site)

        const createEntry = (name, otherProps = {}) => ({ name, ...otherProps })
        const asIndexedList = (entries) => reduce(entries, (acc, entry) => ({
            ...acc,
            [entry.name]: entry
        }), {})

        const createState = (entries) => ({
            sites: {
                entries: asIndexedList(entries),
                otherProp: true
            }
        })

        const queriedEntry = createEntry('proj/abc/sites/x', { primaryDomain: 'x.com' })
        const existingEntries = [
            queriedEntry,
            createEntry('proj/abc/sites/y', { primaryDomain: 'y.com' }),
            createEntry('proj/abc/sites/z', { primaryDomain: 'z.com' })
        ]

        const state = createState(existingEntries)
        const emptyState = createState([])

        describe('getState()', () => {
            it("returns the resource's root state", () => {
                expect(getState(state)).toEqual(state.sites)
                expect(keys(getState(state))).toEqual(['entries', 'otherProp'])
            })

            it('returns the state structure even for empty state', () => {
                expect(getState(emptyState)).toEqual(emptyState.sites)
                expect(keys(getState(emptyState))).toEqual(['entries', 'otherProp'])
            })
        })

        describe('getAll()', () => {
            it('returns all the resource entries in the current state', () => {
                expect(getAll(state)).toEqual(asIndexedList(existingEntries))
            })

            it('returns an empty indexed list for an empty state', () => {
                expect(getAll(emptyState)).toEqual({})
            })
        })

        describe('countAll()', () => {
            it('returns the number of resource entries', () => {
                expect(countAll(state)).toEqual(3)
            })

            it('returns 0 for an empty state', () => {
                expect(countAll(emptyState)).toEqual(0)
            })
        })

        describe('getByName()', () => {
            it('finds an entry by a given name', () => {
                expect(getByName('proj/abc/sites/x')(state)).toEqual(queriedEntry)
            })

            it('returns null if nothing matches', () => {
                expect(getByName('proj/abc/sites/xyz')(state)).toEqual(null)
            })
        })

        describe('getByURL()', () => {
            it('finds an entry by a given URL', () => {
                expect(getForURL('proj/abc/sites/x')(state)).toEqual(queriedEntry)
                expect(getForURL('/proj/abc/sites/x')(state)).toEqual(queriedEntry)
                expect(getForURL('/proj/abc/sites/x/and/some/more?query=true')(state)).toEqual(queriedEntry)
                expect(getForURL('https://dashboard.test/proj/abc/sites/x')(state)).toEqual(queriedEntry)
                expect(getForURL('https://dashboard.test/proj/abc/sites/x?q=test')(state)).toEqual(queriedEntry)
            })

            it('returns null if nothing matches', () => {
                expect(getByName('proj/abc')(state)).toEqual(null)
                expect(getByName('proj/abc/sites')(state)).toEqual(null)
                expect(getByName('proj/abc/sites/xx')(state)).toEqual(null)
                expect(getByName('/orgs/123/proj/abc/sites/xyz')(state)).toEqual(null)
            })
        })
    })

    describe('*saga()', () => {
        const { emitResourceAction } = api
        describe('*emitResourceAction()', () => {
            const actionTypes = createActionTypes(api.Resource.site)

            const createRequest = (method) => ({
                method,
                data: {},
                service: {}
            })

            const createResponse = (method) => ({
                data: {},
                error: null,
                request: createRequest(method)
            })

            describe('emits resource specific actions based on base gRPC action', () => {
                it('for REQUESTED actions', () => {
                    const request = createRequest('listSites')
                    const action = grpc.invoke(request)

                    return expectSaga(emitResourceAction, api.Resource.site, actionTypes, action)
                        .put({
                            type: actionTypes.LIST_REQUESTED,
                            payload: request
                        })
                        .run()
                })

                it('for SUCCEEDED actions', () => {
                    const response = createResponse('listSites')
                    const action = grpc.success(response)

                    return expectSaga(emitResourceAction, api.Resource.site, actionTypes, action)
                        .put({
                            type: actionTypes.LIST_SUCCEEDED,
                            payload: response
                        })
                        .run()
                })

                it('for FAILED actions', () => {
                    const response = createResponse('listSites')
                    const action = grpc.fail(response)

                    return expectSaga(emitResourceAction, api.Resource.site, actionTypes, action)
                        .put({
                            type: actionTypes.LIST_FAILED,
                            payload: response
                        })
                        .run()
                })

                it("ignoring other resource's requests/responses", () => {
                    const response = createResponse('listOrganizations')
                    const action = grpc.success(response)

                    return expectSaga(emitResourceAction, api.Resource.site, actionTypes, action)
                        .run()
                        .then(({ effects }) => expect(size(effects.put)).toEqual(0))
                })
            })
        })
    })

    const { getResourceFromMethodName } = api
    describe('getResourceFromMethodName()', () => {
        it('infers resource type from gRPC method name', () => {
            expect(getResourceFromMethodName('createOrganization')).toEqual(api.Resource.organization)
            expect(getResourceFromMethodName('deleteOrganization')).toEqual(api.Resource.organization)
            expect(getResourceFromMethodName('listProject')).toEqual(api.Resource.project)
            expect(getResourceFromMethodName('callProcedureOnProject')).toEqual(api.Resource.project)
            expect(getResourceFromMethodName('listFoo')).toEqual(null)
            expect(getResourceFromMethodName('invalidRPC')).toEqual(null)
        })
    })

    const { getRequestTypeFromMethodName } = api
    describe('getRequestTypeFromMethodName()', () => {
        it('infers request type from gRPC method name', () => {
            expect(getRequestTypeFromMethodName('createOrganization')).toEqual(api.Request.create)
            expect(getRequestTypeFromMethodName('deleteOrganization')).toEqual(api.Request.destroy)
            expect(getRequestTypeFromMethodName('updateProj')).toEqual(api.Request.update)
            expect(getRequestTypeFromMethodName('listProject')).toEqual(api.Request.list)
            expect(getRequestTypeFromMethodName('listFoo')).toEqual(api.Request.list)
            expect(getRequestTypeFromMethodName('getRes')).toEqual(api.Request.get)
            expect(getRequestTypeFromMethodName('invalidMethod')).toEqual(null)
        })
    })

    const { getStatusFromAction } = api
    describe('getStatusFromAction()', () => {
        const createAction = (type) => ({ type })

        it('infers status from a given gRPC action', () => {
            expect(getStatusFromAction(createAction(grpc.INVOKED))).toEqual(api.Status.requested)
            expect(getStatusFromAction(createAction(grpc.SUCCEEDED))).toEqual(api.Status.succeeded)
            expect(getStatusFromAction(createAction(grpc.FAILED))).toEqual(api.Status.failed)
        })

        it('infers status from a given API action', () => {
            expect(getStatusFromAction(createAction(sites.GET_REQUESTED))).toEqual(api.Status.requested)
            expect(getStatusFromAction(createAction(sites.LIST_SUCCEEDED))).toEqual(api.Status.succeeded)
            expect(getStatusFromAction(createAction(sites.DESTROY_FAILED))).toEqual(api.Status.failed)
        })
    })

    const { getRequestTypeFromAction } = api
    describe('getRequestTypeFromAction()', () => {
        const createRequestActionWithMethod = (type, method) => ({ type, payload: { method } })
        const createResponseActionWithMethod = (type, method) => ({ type, payload: { request : { method } } })

        it('infers status from a given _request_ gRPC action', () => {
            expect(
                getRequestTypeFromAction(createRequestActionWithMethod(grpc.INVOKED, 'listSites'))
            ).toEqual(api.Request.list)

            expect(
                getRequestTypeFromAction(createRequestActionWithMethod(sites.LIST_REQUESTED, 'listSites'))
            ).toEqual(api.Request.list)

            expect(
                getRequestTypeFromAction(createRequestActionWithMethod(grpc.INVOKED, 'noopSites'))
            ).toEqual(null)
        })

        it('infers status from a given _response_ gRPC action', () => {
            expect(
                getRequestTypeFromAction(createResponseActionWithMethod(sites.CREATE_FAILED, 'createSite'))
            ).toEqual(api.Request.create)

            expect(
                getRequestTypeFromAction(createResponseActionWithMethod(sites.UPDATE_FAILED, 'updateOrganizations'))
            ).toEqual(api.Request.update)

            expect(
                getRequestTypeFromAction(createResponseActionWithMethod(grpc.FAILED, 'noopOrganizations'))
            ).toEqual(null)
        })
    })

    describe('resource name helpers', () => {
        const emptyNamePayload = {
            slug   : null,
            name   : null,
            parent : null,
            url    : '/',
            params : {}
        }

        describe('for top-level resource: orgs/:slug', () => {
            const { parseName, buildName } = api.createNameHelpers('orgs/:slug')

            describe('parseName()', () => {
                it('properly parses valid names', () => {
                    expect(parseName('orgs/abc')).toEqual({
                        slug   : 'abc',
                        name   : 'orgs/abc',
                        url    : '/orgs/abc',
                        parent : null,
                        params : {
                            slug: 'abc'
                        }
                    })

                    expect(parseName('orgs/123-abcd')).toEqual({
                        slug   : '123-abcd',
                        name   : 'orgs/123-abcd',
                        url    : '/orgs/123-abcd',
                        parent : null,
                        params : {
                            slug: '123-abcd'
                        }
                    })
                })

                it('properly parses names from full URLs', () => {
                    expect(parseName('/orgs/abc/projects/123/sites/xyz?filter=active')).toEqual({
                        slug   : 'abc',
                        name   : 'orgs/abc',
                        url    : '/orgs/abc',
                        parent : null,
                        params : {
                            slug: 'abc'
                        }
                    })

                    expect(parseName('https://dashboard.test/orgs/abc/projects/123/sites/xyz?filter=active')).toEqual({
                        slug   : 'abc',
                        name   : 'orgs/abc',
                        url    : '/orgs/abc',
                        parent : null,
                        params : {
                            slug: 'abc'
                        }
                    })
                })

                it('returns an empty payload for non-matching names', () => {
                    expect(parseName('orgs/')).toEqual(emptyNamePayload)
                    expect(parseName('proj/abc')).toEqual(emptyNamePayload)
                    expect(parseName('proj/abc/orgs/xyz')).toEqual(emptyNamePayload)
                })
            })

            describe('buildName()', () => {
                it('builds name form given params', () => {
                    expect(buildName({ slug: 'abc' })).toEqual('orgs/abc')
                })

                it('returns null when given invalid/incomplete paylods', () => {
                    expect(buildName({ name: 'test' })).toEqual(null)
                    expect(buildName({})).toEqual(null)
                })
            })
        })

        describe('for nested resource: proj/:proj/sites/:slug', () => {
            const { parseName, buildName } = api.createNameHelpers('proj/:proj/sites/:slug')

            describe('parseName()', () => {
                it('properly parses valid names', () => {
                    expect(parseName('proj/abc/sites/xyz')).toEqual({
                        slug   : 'xyz',
                        name   : 'proj/abc/sites/xyz',
                        url    : '/proj/abc/sites/xyz',
                        parent : 'proj/abc',
                        params : {
                            slug: 'xyz',
                            proj: 'abc'
                        }
                    })
                })

                it('returns an empty payload for non-matching names', () => {
                    expect(parseName('proj/abc')).toEqual(emptyNamePayload)
                    expect(parseName('proj/abc/dev-sites/xyz')).toEqual(emptyNamePayload)
                })
            })

            describe('buildName()', () => {
                it('builds name form given params', () => {
                    expect(buildName({ proj: 'abc', slug: 'xyz' })).toEqual('proj/abc/sites/xyz')
                })

                it('returns null when given invalid/incomplete paylods', () => {
                    expect(buildName({ slug: 'test' })).toEqual(null)
                    expect(buildName({})).toEqual(null)
                })
            })
        })
    })
})
