package source

// SyncAccessControl adds any newly registered APIs and administrator policies
// to an existing database. Each underlying seed is idempotent, so this is safe
// to run during every server start without coupling schema migration to a
// feature-specific permission helper.
func SyncAccessControl() error {
	if err := Api.Init(); err != nil {
		return err
	}
	return Casbin.Init()
}
