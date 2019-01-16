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
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/blaggacao/ridecell-operator/pkg/components"
	instancev1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/instance/v1beta1"
)

type backuperComponent struct {
	templatePath string
}

func NewBackuper(templatePath string) *backuperComponent {
	return &backuperComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=cronjobs/status,verbs=get;update;patch
func (_ *backuperComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&batchv1beta1.CronJob{},
	}
}

func (_ *backuperComponent) IsReconcilable(ctx *components.ComponentContext) bool {
	instance := ctx.Top.(*instancev1beta1.OdooInstance)
	if instance.Spec.ParentHostname == nil {
		// The initializer component is the one that should initialize a root instance
		return false
	}
	if instance.GetStatusCondition(instancev1beta1.OdooInstanceStatusConditionTypeCreated) != nil {
		// The instance is already created (or creating)
		return false
	}
	return true
}

func (comp *backuperComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	instance := ctx.Top.(*instancev1beta1.OdooInstance)

	// Set up the extra data map for the template.
	listObj := &instancev1beta1.OdooInstanceList{}

	res, obj, err := ctx.GetOne(listObj, map[string]string{
		"cluster.odoo.io/name":      instance.Labels["cluster.odoo.io/name"],
		"instance.odoo.io/hostname": *instance.Spec.ParentHostname,
	})
	if err != nil || obj == nil {
		return res, err
	}
	parentInstance := obj.(*instancev1beta1.OdooInstance)
	extra := map[string]interface{}{}
	extra["FromDatabase"] = string(parentInstance.Spec.Hostname)

	obj, err = ctx.GetTemplate(comp.templatePath, extra)
	if err != nil {
		return reconcile.Result{}, err
	}
	cronjob := obj.(*batchv1beta1.CronJob)

	existing := &batchv1beta1.CronJob{}
	err = ctx.Get(ctx.Context, types.NamespacedName{Name: cronjob.Name, Namespace: cronjob.Namespace}, existing)
	if err != nil && errors.IsNotFound(err) {

		instance.SetStatusConditionBackupCronJobSecheduleBackuped()

		// Creating the cronjob
		err = controllerutil.SetControllerReference(instance, cronjob, ctx.Scheme)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = ctx.Create(ctx.Context, cronjob)
		if err != nil {
			// If this fails, someone else might have created a backuper cronjob between the Get and here, so just try again.
			return reconcile.Result{Requeue: true}, err
		}
		// CronJob is created, so we're done for now.
		ctx.Logger.V(1).Info("reconciled", "object", obj, "operation", "created")
		return reconcile.Result{}, nil
	} else if err != nil {
		// Some other real error, bail.
		return reconcile.Result{}, err
	}

	// Job is still running, will get reconciled when it finishes.
	return reconcile.Result{}, nil
}
