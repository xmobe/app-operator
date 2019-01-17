/*
 * This file is part of the Odoo-Operator (R) project.
 * Copyright (c) 2018-2018 XOE Corp. SAS
 * Authors: David Arnold, et al.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * ALTERNATIVE LICENCING OPTION
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving the Odoo-Operator software without
 * disclosing the source code of your own applications. These activities
 * include: Offering paid services to a customer as an ASP, shipping Odoo-
 * Operator with a closed source product.
 *
 */

package components

import (
	"github.com/blaggacao/ridecell-operator/pkg/components"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	instancev1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/instance/v1beta1"
)

type synchMigratorComponent struct {
	templatePath string
}

func NewSyncMigrator(templatePath string) *synchMigratorComponent {
	return &synchMigratorComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
func (_ *synchMigratorComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&batchv1.Job{},
	}
}

func (_ *synchMigratorComponent) IsReconcilable(ctx *components.ComponentContext) bool {
	instance := ctx.Top.(*instancev1beta1.OdooInstance)
	if instance.Spec.ParentHostname == nil {
		// The migrator component is never interfering to initialize a root instance
		return false
	}
	createdCondition := instance.GetStatusCondition(instancev1beta1.OdooInstanceStatusConditionTypeCreated)
	if createdCondition != nil && createdCondition.Status != corev1.ConditionTrue {
		// Apply migrations only on created instances
		return false
	}
	// Get the parent instance ...
	listObj := &instancev1beta1.OdooInstanceList{}
	_, obj, err := ctx.GetOne(listObj, map[string]string{
		"cluster.odoo.io/part-of-cluster": instance.Labels["cluster.odoo.io/part-of-cluster"],
		"app.kubernetes.io/instance":      instance.Labels["instance.odoo.io/part-of-instance"],
	})
	if err != nil || obj == nil {
		return false
	}
	parentInstance := obj.(*instancev1beta1.OdooInstance)

	if instance.Spec.Version == parentInstance.Spec.Version {
		// ... and apply only if there is an explicit version bump over the parent instance
		// TODO: Ensure version order
		return false
	}
	return true
}

func (comp *synchMigratorComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	obj, err := ctx.GetTemplate(comp.templatePath, nil)
	if err != nil {
		return reconcile.Result{}, err
	}
	job := obj.(*batchv1.Job)
	instance := ctx.Top.(*instancev1beta1.OdooInstance)

	existing := &batchv1.Job{}
	err = ctx.Get(ctx.Context, types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, existing)
	if err != nil && errors.IsNotFound(err) {
		instance.SetStatusConditionSynchronousMigratorJobCreationMigrated()

		// Launching the job
		err = controllerutil.SetControllerReference(instance, job, ctx.Scheme)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = ctx.Create(ctx.Context, job)
		if err != nil {
			// If this fails, someone else might have started a copier job between the Get and here, so just try again.
			return reconcile.Result{Requeue: true}, err
		}
		// Job is started, so we're done for now.
		ctx.Logger.V(1).Info("reconciled", "object", obj, "operation", "created")
		return reconcile.Result{}, nil
	} else if err != nil {
		// Some other real error, bail.
		return reconcile.Result{}, err
	}

	// If we get this far, the job previously started at some point and might be done.
	// Check if the job succeeded.
	if existing.Status.Succeeded > 0 {
		// Success! Update the corresponding OdooInstanceStatusCondition and delete the job.

		instance.SetStatusConditionSynchronousMigratorJobSuccessMigrated()
		ctx.Logger.V(1).Info("set", "status contition", "Migrated", "status", "true")

		err = ctx.Delete(ctx.Context, existing, client.PropagationPolicy(metav1.DeletePropagationBackground))
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
		ctx.Logger.V(1).Info("reconciled", "job", existing, "operation", "deleted")
	}

	// ... Or if the job failed.
	if existing.Status.Failed > 0 {
		ctx.Logger.V(1).Info("leaving failed job for debugging", "job", existing)
	}

	// Job is still running, will get reconciled when it finishes.
	return reconcile.Result{}, nil
}
