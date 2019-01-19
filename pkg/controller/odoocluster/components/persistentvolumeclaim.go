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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/blaggacao/ridecell-operator/pkg/components"
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

func (comp *persistentVolumeClaimComponent) IsReconcilable(ctx *components.ComponentContext) bool {
	// PersistentVaolumeClaims are immutable after creation
	// TODO: Forbidden: is immutable after creation except resources.requests for bound claims"
	obj, err := ctx.GetTemplate(comp.templatePath, nil)
	if err != nil {
		return false
	}
	fetchObj := obj.(*corev1.PersistentVolumeClaim)

	err = ctx.Get(ctx.Context, types.NamespacedName{Name: fetchObj.Name, Namespace: fetchObj.Namespace}, fetchObj)
	if err != nil && errors.IsNotFound(err) {
		return true
	}
	return false
}

func (comp *persistentVolumeClaimComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	res, _, err := ctx.CreateOrUpdate(comp.templatePath, nil, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*corev1.PersistentVolumeClaim)
		existing := existingObj.(*corev1.PersistentVolumeClaim)
		// Copy the Spec over.
		if &existing.Spec != nil {
			existing.Spec = goal.Spec
		} else {
			existing.Spec = corev1.PersistentVolumeClaimSpec{}
		}
		return nil
	})
	return res, err
}
