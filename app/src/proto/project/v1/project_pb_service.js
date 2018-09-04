// package: 
// file: project/v1/project.proto

var project_v1_project_pb = require("../../project/v1/project_pb");
var grpc = require("grpc-web-client").grpc;

var Projects = (function () {
  function Projects() {}
  Projects.serviceName = "Projects";
  return Projects;
}());

Projects.List = {
  methodName: "List",
  service: Projects,
  requestStream: false,
  responseStream: true,
  requestType: project_v1_project_pb.ListRequest,
  responseType: project_v1_project_pb.Project
};

exports.Projects = Projects;

function ProjectsClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ProjectsClient.prototype.list = function list(requestMessage, metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.invoke(Projects.List, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onMessage: function (responseMessage) {
      listeners.data.forEach(function (handler) {
        handler(responseMessage);
      });
    },
    onEnd: function (status, statusMessage, trailers) {
      listeners.end.forEach(function (handler) {
        handler();
      });
      listeners.status.forEach(function (handler) {
        handler({ code: status, details: statusMessage, metadata: trailers });
      });
      listeners = null;
    }
  });
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

exports.ProjectsClient = ProjectsClient;

