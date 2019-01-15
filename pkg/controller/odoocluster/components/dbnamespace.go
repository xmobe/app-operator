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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type dbNamespaceComponent struct {
	templatePath string
}

func NewDBNamespace(templatePath string) *dbNamespaceComponent {
	return &dbNamespaceComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=appscluster.odoo.io,resources=dbnamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=appscluster.odoo.io,resources=dbnamespaces/status,verbs=get;update;patch
func (_ *dbNamespaceComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&clusterv1beta1.DBNamespace{},
	}
}

func (_ *dbNamespaceComponent) IsReconcilable(_ *components.ComponentContext) bool {
	// DBNamespaces have no dependencies, always reconcile.
	return true
}

func (comp *dbNamespaceComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {

	res, _, err := ctx.CreateOrUpdate(comp.templatePath, nil, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*clusterv1beta1.DBNamespace)
		existing := existingObj.(*clusterv1beta1.DBNamespace)
		// Copy the Spec over.
		existing.Spec = goal.Spec
		return nil
	})
	return res, err
}
