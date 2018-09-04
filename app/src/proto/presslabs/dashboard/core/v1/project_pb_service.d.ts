// package: presslabs.dashboard.v1
// file: presslabs/dashboard/core/v1/project.proto

import * as presslabs_dashboard_core_v1_project_pb from "../../../../presslabs/dashboard/core/v1/project_pb";
import {grpc} from "grpc-web-client";

type ProjectsList = {
  readonly methodName: string;
  readonly service: typeof Projects;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof presslabs_dashboard_core_v1_project_pb.ListRequest;
  readonly responseType: typeof presslabs_dashboard_core_v1_project_pb.Project;
};

export class Projects {
  static readonly serviceName: string;
  static readonly List: ProjectsList;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }
export type ServiceClientOptions = { transport: grpc.TransportConstructor; debug?: boolean }

interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: () => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}

export class ProjectsClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: ServiceClientOptions);
  list(requestMessage: presslabs_dashboard_core_v1_project_pb.ListRequest, metadata?: grpc.Metadata): ResponseStream<presslabs_dashboard_core_v1_project_pb.Project>;
}

