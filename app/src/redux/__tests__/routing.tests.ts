import { routing } from '../'
import { createLocation } from 'history'

describe('api', () => {
    describe('matchRoute()', () => {
        const { matchRoute } = routing

        it('matches simple routes', () => {
            const route = createLocation({ path: 'dashboard' })
            expect(matchRoute(route)).toEqual({
                key     : 'dashboard',
                isExact : true,
                params  : {},
                path    : '/',
                url     : '/'
            })
        })

        it('matches routes with query arguments', () => {
            const route = createLocation({ pathname: '/', search: '?filter=active' })
            expect(matchRoute(route)).toEqual({
                key     : 'dashboard',
                isExact : true,
                params  : { filter: 'active' },
                path    : '/',
                url     : '/?filter=active'
            })
        })

        it('matches routes with path segments', () => {
            const route = createLocation({ pathname: '/orgs/test', search: '?filter=active' })
            expect(matchRoute(route)).toEqual({
                key     : 'organizations',
                isExact : true,
                params  : { slug: 'test', filter: 'active', action: undefined },
                path    : '/orgs/:slug?/:action?',
                url     : '/orgs/test?filter=active'
            })
        })
    })

    describe('routeFor()', () => {
        const { routeFor } = routing

        it('builds URLs based on the routing table', () => {
            expect(routeFor('dashboard')).toEqual('/')
        })

        it('replaces path segments with given params', () => {
            expect(routeFor('organizations', { slug: 'test' })).toEqual('/orgs/test')
            expect(routeFor('organizations', { slug: '_', action: 'new' })).toEqual('/orgs/_/new')
            expect(routeFor('organizations', { action: 'new' })).toEqual('/orgs/new')
        })

        it('appends extra parameters to the query', () => {
            expect(routeFor('organizations', { slug: 'test', filter: 'active' })).toEqual('/orgs/test?filter=active')
        })

        it('throws if a required path segment is missing', () => {
            expect(() => { routeFor('sites') }).toThrow(new TypeError('Expected "project" to be a string'))
        })

        it('throws if an invalid key is used', () => {
            expect(() => { routeFor('someKey') }).toThrow(new Error('Invalid route key: someKey'))
        })
    })


    describe('routeForResource()', () => {
        const { routeForResource } = routing

        it("builds URLs using a resource's name", () => {
            const resource = { name: 'orgs/test' }
            expect(routeForResource(resource)).toEqual('/orgs/test')
        })

        it('fills all path segments properly', () => {
            const resource = { name: 'orgs/test' }
            expect(routeForResource(resource, { action: 'edit' })).toEqual('/orgs/test/edit')
        })

        it('appends query arguments accordingly', () => {
            const resource = { name: 'orgs/test' }
            expect(routeForResource(resource, { filter: 'active' })).toEqual('/orgs/test?filter=active')
        })

        it('works with non top-level resources', () => {
            const resource = { name: 'project/123/sites/abc' }
            expect(routeForResource(resource, { filter: 'active' })).toEqual('/project/123/sites/abc?filter=active')
        })

        it('works with any type of resource', () => {
            const resource = { name: 'resource/xxx' }
            expect(routeForResource(resource, { filter: 'active' })).toEqual('/resource/xxx?filter=active')
        })
    })
})
