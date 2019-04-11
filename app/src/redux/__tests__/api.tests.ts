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

                it('properly parses names from full (maybe longer) URLs', () => {
                    expect(parseName('/orgs/abc/projects/123/sites/xyz?filter=active')).toEqual({
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
