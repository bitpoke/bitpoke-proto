// package: 
// file: projects/v1/project.proto

import * as projects_v1_project_pb from "../../projects/v1/project_pb";
import {grpc} from "grpc-web-client";

type ProjectsListProjects = {
  readonly methodName: string;
  readonly service: typeof Projects;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof projects_v1_project_pb.ListProjectsRequest;
  readonly responseType: typeof projects_v1_project_pb.Project;
};

export class Projects {
  static readonly serviceName: string;
  static readonly ListProjects: ProjectsListProjects;
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
  listProjects(requestMessage: projects_v1_project_pb.ListProjectsRequest, metadata?: grpc.Metadata): ResponseStream<projects_v1_project_pb.Project>;
}

