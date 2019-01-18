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
)

type serviceComponent struct {
	templatePath string
}

func NewService(templatePath string) *serviceComponent {
	return &serviceComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=services/status,verbs=get;update;patch
func (_ *serviceComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		// Todo own OdooInstances by the cluster on creation
		&corev1.Service{},
	}
}

func (_ *serviceComponent) IsReconcilable(_ *components.ComponentContext) bool {
	return true
}

func (comp *serviceComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	// Set up the extra data map for the template.
	extra := map[string]interface{}{}
	extra["Service"] = true

	res, _, err := ctx.CreateOrUpdate(comp.templatePath, extra, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*corev1.Service)
		existing := existingObj.(*corev1.Service)
		// Copy the configuration Data over.
		existingClusterIP := existing.Spec.ClusterIP
		existing.Spec = goal.Spec
		// ClusterIP is immutable after creation
		if !existing.ObjectMeta.CreationTimestamp.IsZero() {
			existing.Spec.ClusterIP = existingClusterIP
		}
		return nil
	})
	return res, err
}
