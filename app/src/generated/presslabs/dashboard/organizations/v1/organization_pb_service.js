// package: presslabs.dashboard.organizations.v1
// file: presslabs/dashboard/organizations/v1/organization.proto

var presslabs_dashboard_organizations_v1_organization_pb = require("../../../../presslabs/dashboard/organizations/v1/organization_pb");
var google_protobuf_empty_pb = require("google-protobuf/google/protobuf/empty_pb");
var grpc = require("grpc-web-client").grpc;

var OrganizationsService = (function () {
  function OrganizationsService() {}
  OrganizationsService.serviceName = "presslabs.dashboard.organizations.v1.OrganizationsService";
  return OrganizationsService;
}());

OrganizationsService.CreateOrganization = {
  methodName: "CreateOrganization",
  service: OrganizationsService,
  requestStream: false,
  responseStream: false,
  requestType: presslabs_dashboard_organizations_v1_organization_pb.CreateOrganizationRequest,
  responseType: presslabs_dashboard_organizations_v1_organization_pb.Organization
};

OrganizationsService.GetOrganization = {
  methodName: "GetOrganization",
  service: OrganizationsService,
  requestStream: false,
  responseStream: false,
  requestType: presslabs_dashboard_organizations_v1_organization_pb.GetOrganizationRequest,
  responseType: presslabs_dashboard_organizations_v1_organization_pb.Organization
};

OrganizationsService.UpdateOrganization = {
  methodName: "UpdateOrganization",
  service: OrganizationsService,
  requestStream: false,
  responseStream: false,
  requestType: presslabs_dashboard_organizations_v1_organization_pb.UpdateOrganizationRequest,
  responseType: presslabs_dashboard_organizations_v1_organization_pb.Organization
};

OrganizationsService.DeleteOrganization = {
  methodName: "DeleteOrganization",
  service: OrganizationsService,
  requestStream: false,
  responseStream: false,
  requestType: presslabs_dashboard_organizations_v1_organization_pb.DeleteOrganizationRequest,
  responseType: google_protobuf_empty_pb.Empty
};

OrganizationsService.ListOrganizations = {
  methodName: "ListOrganizations",
  service: OrganizationsService,
  requestStream: false,
  responseStream: false,
  requestType: presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsRequest,
  responseType: presslabs_dashboard_organizations_v1_organization_pb.ListOrganizationsResponse
};

exports.OrganizationsService = OrganizationsService;

function OrganizationsServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

OrganizationsServiceClient.prototype.createOrganization = function createOrganization(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationsService.CreateOrganization, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

OrganizationsServiceClient.prototype.getOrganization = function getOrganization(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationsService.GetOrganization, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

OrganizationsServiceClient.prototype.updateOrganization = function updateOrganization(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationsService.UpdateOrganization, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

OrganizationsServiceClient.prototype.deleteOrganization = function deleteOrganization(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationsService.DeleteOrganization, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

OrganizationsServiceClient.prototype.listOrganizations = function listOrganizations(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(OrganizationsService.ListOrganizations, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.OrganizationsServiceClient = OrganizationsServiceClient;

