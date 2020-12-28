
# VERIFY TENANT PATCH
## GOAL
This policy verifies that a tenant patch is allowed. In particular: a basic user can update only the publicKeys fields while an admin can update anything,
## TESTS
Tests are available in folder [policies](./policies).
## HOW TO DEPLOY
Run the following commands inside folder [manifest](./manifest):
- kubectl create -f config_sync.yaml
- kubectl create -f template.yaml
- kubectl create -f constraint.yaml

**Severity:** Violation

**Resources:** [Tenant](../../operators/deploy/crds/crownlabs.polito.it_tenants.yaml)



