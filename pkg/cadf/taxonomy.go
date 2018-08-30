package cadf

// IsValidTypeURI that matches CADF Taxonomy.
// Was a nice idea, but I need to check if each one starts with...
func IsValidTypeURI(TypeURI string) bool {
    switch TypeURI {
    case
		"storage",
		"storage/node",
		"storage/volume",
		"storage/memory",
		"storage/container",
		"storage/directory",
		"storage/database",
		"storage/queue",
		"compute",
		"compute/node",
		"compute/cpu",
		"compute/machine",
		"compute/process",
		"compute/thread",
		"network",
		"network/node",
		"network/node/host",
		"network/connection",
		"network/domain",
		"network/cluster",
		"service",
		"service/oss",
		"service/bss",
		"service/bss/metering",
		"service/composition",
		"service/compute",
		"service/database",
		"service/security",
		"service/security/keymanager",
		"service/security/account",
		"service/security/account/user",
		"service/security/audit/filter",
		"service/storage",
		"service/storage/block",
		"service/storage/image",
		"service/storage/object",
		"service/network",
		"data",
		"data/message",
		"data/workload",
		"data/workload/app",
		"data/workload/service",
		"data/workload/task",
		"data/workload/job",
		"data/file",
		"data/file/catalog",
		"data/file/log",
		"data/template",
		"data/package",
		"data/image",
		"data/module",
		"data/config",
		"data/directory",
		"data/database",
		"data/security",
		"data/security/account",
		"data/security/credential",
		"data/security/domain",
		"data/security/endpoint",
		"data/security/group",
		"data/security/identity",
		"data/security/key",
		"data/security/license",
		"data/security/policy",
		"data/security/profile",
		"data/security/project",
		"data/security/region",
		"data/security/role",
		"data/security/service",
		"data/security/trust",
		"data/security/account/user",
		"data/security/account/user/privilege",
		"data/database/alias",
		"data/database/catalog",
		"data/database/constraints",
		"data/database/index",
		"data/database/instance",
		"data/database/key",
		"data/database/routine",
		"data/database/schema",
		"data/database/sequence",
		"data/database/table",
		"data/database/trigger",
		"data/database/view":
        return true
    }
    return false
}