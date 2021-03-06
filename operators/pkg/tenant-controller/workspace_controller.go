/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tenant_controller

import (
	"context"
	"fmt"

	crownlabsv1alpha1 "github.com/netgroup-polito/CrownLabs/operators/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WorkspaceReconciler reconciles a Workspace object
type WorkspaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	KcA    *KcActor
}

// +kubebuilder:rbac:groups=crownlabs.polito.it,resources=workspaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=crownlabs.polito.it,resources=workspaces/status,verbs=get;update;patch

// Reconcile reconciles the state of a workspace resource
func (r *WorkspaceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	var ws crownlabsv1alpha1.Workspace

	if err := r.Get(ctx, req.NamespacedName, &ws); err != nil {
		// reconcile was triggered by a delete request
		klog.Infof("Workspace %s deleted", req.Name)
		rolesToDelete := genWorkspaceRoleNames(req.Name)
		if err := r.KcA.deleteKcRoles(ctx, rolesToDelete); err != nil {
			klog.Errorf("Error when deleting roles of workspace %s", req.NamespacedName)
			klog.Error(err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var retrigErr error = nil
	if ws.Status.Subscriptions == nil {
		ws.Status.Subscriptions = make(map[string]crownlabsv1alpha1.SubscriptionStatus)
	}
	klog.Infof("Reconciling workspace %s", req.Name)

	nsName := fmt.Sprintf("workspace-%s", ws.Name)
	ns := v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: nsName}}

	nsOpRes, err := ctrl.CreateOrUpdate(ctx, r.Client, &ns, func() error {
		updateNamespace(&ns)
		return ctrl.SetControllerReference(&ws, &ns, r.Scheme)
	})
	if err != nil {
		klog.Errorf("Unable to create or update namespace of workspace %s", ws.Name)
		klog.Error(err)
		ws.Status.Namespace.Created = false
		ws.Status.Namespace.Name = ""
		retrigErr = err
	} else {
		klog.Infof("Namespace %s for workspace %s %s", nsName, req.Name, nsOpRes)
		ws.Status.Namespace.Created = true
		ws.Status.Namespace.Name = nsName
	}

	if err := r.KcA.createKcRoles(ctx, genWorkspaceRoleNames(ws.Name)); err != nil {
		ws.Status.Subscriptions["keycloak"] = crownlabsv1alpha1.SubscrFailed
		retrigErr = err
	} else {
		ws.Status.Subscriptions["keycloak"] = crownlabsv1alpha1.SubscrOk
	}

	// everything should went ok, update status before exiting reconcile
	if err := r.Status().Update(ctx, &ws); err != nil {
		// if status update fails, still try to reconcile later
		klog.Error("Unable to update status before exiting reconciler", err)
		retrigErr = err
	}

	return ctrl.Result{}, retrigErr
}

func (r *WorkspaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crownlabsv1alpha1.Workspace{}).
		Complete(r)
}

func updateNamespace(ns *v1.Namespace) {
	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}
	ns.Labels["crownlabs.polito.it/type"] = "workspace"
}

func genWorkspaceRoleNames(wsName string) []string {
	return []string{fmt.Sprintf("workspace-%s:user", wsName), fmt.Sprintf("workspace-%s:admin", wsName)}
}
