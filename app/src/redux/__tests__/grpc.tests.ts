import { expectSaga } from 'redux-saga-test-plan'
import { presslabs } from '@presslabs/dashboard-proto'

import { omitBy, isNil } from 'lodash'

import { grpc, sites } from '../'

const { SitesService } = presslabs.dashboard.sites.v1

describe('gRPC', () => {
    describe('*saga()', () => {

        describe('*performRequest()', () => {

            const createFakeTransport = (error = null, responseData = {}) => {
                return (method, requestData, callback) => {
                    if (error) {
                        callback(error, null)
                        return
                    }

                    callback(null, responseData)
                }
            }

            const createService = (error, responseData) =>
                SitesService.create(createFakeTransport(error, responseData))

            const createRequest = (service) => ({
                method: 'listSites',
                data: [],
                service
            })

            const createResponse = (request, error, data) => (omitBy({
                request,
                error,
                data: data ? new sites.ListSitesResponse(data) : null
            }, isNil))

            it('handles requests that complete successfully', () => {
                const responseData = []
                const responseError = null

                const service = createService(responseError, responseData)
                const request = createRequest(service)
                const response = createResponse(request, responseError, responseData)

                return expectSaga(grpc.saga)
                    .dispatch(grpc.invoke(request))
                    .put(grpc.success(response))
                    .silentRun()
            })

            it('handles requests that fail', () => {
                const responseData = null
                const responseError = new Error()

                const service = createService(responseError, responseData)
                const request = createRequest(service)
                const response = createResponse(request, responseError, responseData)

                return expectSaga(grpc.saga)
                    .dispatch(grpc.invoke(request))
                    .put(grpc.fail(response))
                    .silentRun()
            })
        })
    })
})
