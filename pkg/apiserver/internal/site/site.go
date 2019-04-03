package site

import (
	"github.com/presslabs/dashboard/pkg/apiserver/status"
	"github.com/presslabs/dashboard/pkg/internal/projectns"
	"k8s.io/apimachinery/pkg/api/errors"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// Resolve resolves an fully-qualified site name to a k8s object name.
// The function returns site name, project name and error
func Resolve(name string) (string, string, error) {
	if path.Clean(name) != name {
		return "", "", status.InvalidArgumentf("site resources fully-qualified name must be in form project/PROJECT-NAME/site/SITE-NAME")
	}

	matched, err := path.Match("project/*/site/*", name)
	if err != nil || !matched {
		return "", "", status.InvalidArgumentf("site resources fully-qualified name must be in form project/PROJECT-NAME/site/SITE-NAME")
	}

	names := strings.Split(name, "/")
	return names[3], names[1], nil
}

// ResolveToObjectKey resolves an fully-qualified site name to a k8s object name.
// The function returns the object key from FQName and an error
func ResolveToObjectKey(c client.Client, fqSiteName, orgName string) (*client.ObjectKey, error) {
	siteName, projName, err := Resolve(fqSiteName)
	if err != nil {
		return nil, err
	}

	ns, err := projectns.Lookup(c, projName, orgName)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, status.NotFoundf("project not found")
		}
		return nil, status.InternalError()
	}

	key := client.ObjectKey{
		Name:      siteName,
		Namespace: ns.Name,
	}
	return &key, nil
}
