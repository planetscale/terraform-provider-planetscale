package hooks

/*
 * This file is only ever generated once on the first generation and then is free to be modified.
 * Any hooks you wish to add should be registered in the initHooks function. Feel free to define
 * your hooks in this file or in separate files in the hooks package.
 *
 * Hooks are registered per SDK instance, and are valid for the lifetime of the SDK instance.
 */

func initHooks(h *Hooks) {
	customSecurityHook := &CustomSecurityHook{}
	nekiShardReassignmentHook := NewNekiShardReassignmentHook()

	h.registerSDKInitHook(NewPostgresBranchNoContentSkipHook())
	h.registerSDKInitHook(NewPostgresBranchRegionSlugHook())
	h.registerSDKInitHook(NewPostgresBouncerNoContentSkipHook())
	h.registerSDKInitHook(NewVitessBranchNoContentSkipHook())
	h.registerSDKInitHook(NewNekiParametersHook())
	h.registerSDKInitHook(&NekiExtensionsHook{})
	h.registerSDKInitHook(NewReadOnlyReplicaRegionSlugHook())
	h.registerBeforeRequestHook(customSecurityHook)
	h.registerBeforeRequestHook(nekiShardReassignmentHook)
	h.registerAfterSuccessHook(nekiShardReassignmentHook)
	h.registerAfterSuccessHook(NewClientErrorHook())
	// h.registerAfterErrorHook(exampleHook)
	// h.registerAfterSuccessHook(exampleHook)
}
