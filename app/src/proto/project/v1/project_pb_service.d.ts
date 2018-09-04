// package: 
// file: project/v1/project.proto

import * as project_v1_project_pb from "../../project/v1/project_pb";
import {grpc} from "grpc-web-client";

type ProjectsList = {
  readonly methodName: string;
  readonly service: typeof Projects;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof project_v1_project_pb.ListRequest;
  readonly responseType: typeof project_v1_project_pb.Project;
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
  list(requestMessage: project_v1_project_pb.ListRequest, metadata?: grpc.Metadata): ResponseStream<project_v1_project_pb.Project>;
}

