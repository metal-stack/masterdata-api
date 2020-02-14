package auth

import "github.com/metal-stack/security"
import "github.com/metal-stack/metal-lib/jwt/sec"

var (
	// Edit Groupname
	EditGroups = []security.ResourceAccess{
		security.ResourceAccess("tmdm-all-all-edit"),
	}

	EditAccess = sec.MergeResourceAccess(EditGroups)

	// Edit User
	EditUser = security.User{
		EMail:  "tmdm@metal-stack.io",
		Name:   "tmdm",
		Groups: sec.MergeResourceAccess(EditGroups),
		Tenant: "tmdm",
	}
)
