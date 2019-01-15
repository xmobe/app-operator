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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type persistentVolumeClaimComponent struct {
	templatePath string
}

func NewPersistentVolumeClaim(templatePath string) *persistentVolumeClaimComponent {
	return &persistentVolumeClaimComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=persistentvolumeclaims/status,verbs=get;update;patch
func (_ *persistentVolumeClaimComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&corev1.PersistentVolumeClaim{},
	}
}

func (_ *persistentVolumeClaimComponent) IsReconcilable(ctx *components.ComponentContext) bool {
	instance := ctx.Top.(*clusterv1beta1.OdooCluster)
	if !instance.ObjectMeta.CreationTimestamp.IsZero() {
		// PersistentVaolumeClaims are immutable after creation
		return false
	}
	return true
}

func (comp *persistentVolumeClaimComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	extra := map[string]interface{}{}
	extra["VolumeSpec"] = "ok"
	res, _, err := ctx.CreateOrUpdate(comp.templatePath, extra, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*corev1.ConfigMap)
		existing := existingObj.(*corev1.ConfigMap)
		// Copy the configuration Data over.
		existing.Data = goal.Data
		return nil
	})
	return res, err
}
