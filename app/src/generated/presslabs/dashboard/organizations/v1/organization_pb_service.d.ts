// package: presslabs.dashboard.organizations.v1
// file: presslabs/dashboard/organizations/v1/organization.proto

import * as presslabs_dashboard_organizations_v1_organization_pb from "../../../../presslabs/dashboard/organizations/v1/organization_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import {grpc} from "grpc-web-client";

type OrganizationsServiceCreateOrganization = {
  readonly methodName: string;
  readonly service: typeof OrganizationsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof presslabs_dashboard_organizations_v1_organization_pb.CreateOrganizationRequest;
  readonly responseType: typeof presslabs_dashboard_organizations_v1_organization_pb.Organization;
};

type OrganizationsServiceGetOrganization = {
  readonly methodName: string;
  readonly service: typeof OrganizationsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof presslabs_dashboard_organizations_v1_organization_pb.GetOrganizationRequest;
  readonly responseType: typeof presslabs_dashboard_organizations_v1_organization_pb.Organization;
};

type OrganizationsServiceUpdateOrganization = {
  readonly methodName: string;
  readonly service: typeof OrganizationsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof presslabs_dashboard_organizations_v1_organization_pb.UpdateOrganizationRequest;
  readonly responseType: typeof presslabs_dashboard_organizations_v1_organization_pb.Organization;
};

type OrganizationsServiceDeleteOrganization = {
  readonly methodName: string;
  readonly service: typeof OrganizationsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof presslabs_dashboard_organizations_v1_organization_pb.DeleteOrganizationRequest;
  readonly responseType: typeof google_protobuf_empty_pb.Empty;
};

type OrganizationsServiceListOrganizations = {
  readonly methodName: string;
  readonly service: typeof OrganizationsService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsRequest;
  readonly responseType: typeof presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsResponse;
};

export class OrganizationsService {
  static readonly serviceName: string;
  static readonly CreateOrganization: OrganizationsServiceCreateOrganization;
  static readonly GetOrganization: OrganizationsServiceGetOrganization;
  static readonly UpdateOrganization: OrganizationsServiceUpdateOrganization;
  static readonly DeleteOrganization: OrganizationsServiceDeleteOrganization;
  static readonly ListOrganizations: OrganizationsServiceListOrganizations;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: () => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: () => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: () => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class OrganizationsServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  createOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.CreateOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.Organization|null) => void
  ): UnaryResponse;
  createOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.CreateOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.Organization|null) => void
  ): UnaryResponse;
  getOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.GetOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.Organization|null) => void
  ): UnaryResponse;
  getOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.GetOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.Organization|null) => void
  ): UnaryResponse;
  updateOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.UpdateOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.Organization|null) => void
  ): UnaryResponse;
  updateOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.UpdateOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.Organization|null) => void
  ): UnaryResponse;
  deleteOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.DeleteOrganizationRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  deleteOrganization(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.DeleteOrganizationRequest,
    callback: (error: ServiceError|null, responseMessage: google_protobuf_empty_pb.Empty|null) => void
  ): UnaryResponse;
  listOrganizations(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsResponse|null) => void
  ): UnaryResponse;
  listOrganizations(
    requestMessage: presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsRequest,
    callback: (error: ServiceError|null, responseMessage: presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsResponse|null) => void
  ): UnaryResponse;
}

