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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	instancev1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/instance/v1beta1"
	"github.com/xoe-labs/odoo-operator/pkg/labels"
)

type metaComponent struct {
}

func NewMeta() *metaComponent {
	return &metaComponent{}
}

func (_ *metaComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{}
}

func (_ *metaComponent) IsReconcilable(_ *components.ComponentContext) bool {
	// Always "reconcile" (apply) Meta as the second controller action
	return true
}

func (comp *metaComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {

	instance := ctx.Top.(*instancev1beta1.OdooInstance)
	err := labels.AddLabels(ctx.Top, map[string]string{
		"cluster.odoo.io/name":         instance.Spec.Cluster,
		"cluster.odoo.io/track":        "",
		"instance.odoo.io/hostname":    instance.Spec.Hostname,
		"app.kubernetes.io/name":       "odooinstance",
		"app.kubernetes.io/instance":   instance.Name,
		"app.kubernetes.io/component":  "app",
		"app.kubernetes.io/managed-by": "odoo-operator",
		"app.kubernetes.io/version":    instance.Spec.Version,
	})

	op, err := controllerutil.CreateOrUpdate(ctx.Context, ctx, instance.DeepCopyObject(), func(existing runtime.Object) error {
		// Sync the metadata fields.
		targetMeta := ctx.Top.(metav1.ObjectMetaAccessor).GetObjectMeta().(*metav1.ObjectMeta)
		existingMeta := existing.(metav1.ObjectMetaAccessor).GetObjectMeta().(*metav1.ObjectMeta)
		return components.ReconcileMeta(targetMeta, existingMeta)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	ctx.Logger.V(1).Info("reconciled", "odoo instance", instance, "operation", op)
	return reconcile.Result{}, nil
}
