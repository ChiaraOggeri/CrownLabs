package crownlabs_tenant_patch

#check if creator has a tenant 
violation[{"msg": msg, "details": {}}] {
	user := input.review.userInfo.username
	not data.inventory.cluster["crownlabs.polito.it/v1alpha1"].Tenant[user]
	msg := sprintf("Denied request.There is no tenant resource for user %v in the cluster", [user])
}

#check if creator has already a tenant
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "CREATE"
	user := input.review.userInfo.username
	tenantName := input.review.object.metadata.name
	user == tenantName
	msg := sprintf("Denied creation request. %v already has a tenant", [user])
}

#check creation of new user can be done only in workspaces where the creatore is enrolled
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "CREATE"
	user := input.review.userInfo.username
	requiredWs := {ws | ws := input.review.object.spec.workspaces[_].workspaceRef.name}
	providedWs := {ws | ws := data.inventory.cluster["crownlabs.polito.it/v1alpha1"].Tenant[user].spec.workspaces[_].workspaceRef.name}
	missing := requiredWs - providedWs
	count(missing) > 0
	msg := sprintf("Denied creation request. You don't have the priviledges to create users in this workspace: %v", [missing])
}

#check creation of new user can be done only in workspaces where the creator is admin
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "CREATE"
	user := input.review.userInfo.username
	some p
	some q
	ws := input.review.object.spec.workspaces[p].workspaceRef.name
	wsData := data.inventory.cluster["crownlabs.polito.it/v1alpha1"].Tenant[user].spec.workspaces[q]
	wsData.workspaceRef.name == ws
	wsData.role != "admin"
	msg := sprintf("Denied creation request. You don't have the priviledges to create users in this workspace: %v", [ws])
}

#update can't be done on anything but publicKeys and workspaces
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "UPDATE"
	user := input.review.userInfo.username
	some key
	input.review.oldObject.spec[key] != input.review.object.spec[key]
	key != "publicKeys"
	key != "workspaces"
	msg := sprintf("Denied patch request. You can't modify: %v field of tenant resource", [key])
}

#update can't be done on a workspace the creator is not part of
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "UPDATE"
	user := input.review.userInfo.username
	some key
	input.review.oldObject.spec[key] != input.review.object.spec[key]
	key == "workspaces"
	some i
	input.review.oldObject.spec.workspaces[i] != input.review.object.spec.workspaces[i]
	requiredWs := {ws | ws := input.review.object.spec.workspaces[i].workspaceRef.name}
	providedWs := {ws | ws := data.inventory.cluster["crownlabs.polito.it/v1alpha1"].Tenant[user].spec.workspaces[_].workspaceRef.name}
	missing := requiredWs - providedWs
	count(missing) > 0

	msg := sprintf("Denied patch request. You can't modify the workspace %v, because you are not enrolled in it", [missing])
}

#update can't be done on a workspace you are not admin for
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "UPDATE"
	user := input.review.userInfo.username
	some key
	input.review.oldObject.spec[key] != input.review.object.spec[key]
	key == "workspaces"
	some i
	input.review.oldObject.spec.workspaces[i] != input.review.object.spec.workspaces[i]
	wsChangedName := input.review.object.spec.workspaces[i].workspaceRef.name
	some j
	wsRole = data.inventory.cluster["crownlabs.polito.it/v1alpha1"].Tenant[user].spec.workspaces[j].role
	wsName = data.inventory.cluster["crownlabs.polito.it/v1alpha1"].Tenant[user].spec.workspaces[j].workspaceRef.name
	wsName == wsChangedName
	wsRole != "admin"
	msg := sprintf("Denied patch request. You have not the priviledges to modify %v workspace", [wsName])
}

#update on publicKeys can be done only on your own tenant
violation[{"msg": msg, "details": {}}] {
	operation := input.review.operation
	operation == "UPDATE"
	user := input.review.userInfo.username
	tenantName := input.review.oldObject.metadata.name
	oldObj := input.review.oldObject
	newObj := input.review.object
	some key
	oldObj.spec[key] != newObj.spec[key]
	key == "publicKeys"
	user != tenantName
	msg := sprintf("Denied patch request. You don't have the priviledges to modify: %v field", [key])
}
