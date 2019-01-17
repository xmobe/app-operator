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
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
	instancev1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/instance/v1beta1"
	// odooinstanceutils "github.com/xoe-labs/odoo-operator/pkg/controller/odooinstance/utils"
)

type ingressComponent struct {
	templatePath string
}

func NewIngress(templatePath string) *ingressComponent {
	return &ingressComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=extensions,resources=ingress,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=ingress/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=instance.odoo.io,resources=odooinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=instance.odoo.io,resources=odooinstances/status,verbs=get;update;patch
func (_ *ingressComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&instancev1beta1.OdooInstance{},
		&extensionsv1beta1.Ingress{},
	}
}

func (_ *ingressComponent) WatchPredicateFuncs() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(_ event.UpdateEvent) bool {
			return true
		},
		CreateFunc: func(_ event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(_ event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(_ event.GenericEvent) bool {
			return false
		},
	}
}

func (_ *ingressComponent) IsReconcilable(ctx *components.ComponentContext) bool {
	return true
}

func (comp *ingressComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	instance := ctx.Top.(*clusterv1beta1.OdooCluster)

	// Fetch OdooInstances to populate rules
	odooInstanceList := &instancev1beta1.OdooInstanceList{}
	listoptions := client.InNamespace(instance.Namespace)
	listoptions.MatchingLabels(map[string]string{
		"cluster.odoo.io/part-of-cluster": instance.Name,
	})
	err := ctx.List(ctx.Context, listoptions, odooInstanceList)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Prepare extra data with OdooInstanceList
	extra := map[string]interface{}{}
	extra["InstanceList"] = odooInstanceList
	res, _, err := ctx.CreateOrUpdate(comp.templatePath, extra, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*extensionsv1beta1.Ingress)
		existing := existingObj.(*extensionsv1beta1.Ingress)
		// Copy the Spec over.
		existing.Spec = goal.Spec
		return nil
	})
	return res, err
}
