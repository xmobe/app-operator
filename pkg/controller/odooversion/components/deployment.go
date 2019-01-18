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
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type deploymentComponent struct {
	templatePath string
}

func NewDeployment(templatePath string) *deploymentComponent {
	return &deploymentComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
func (_ *deploymentComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&appsv1.Deployment{},
	}
}

func (_ *deploymentComponent) IsReconcilable(_ *components.ComponentContext) bool {
	return true
}

func (comp *deploymentComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	instance := ctx.Top.(*clusterv1beta1.OdooVersion)
	extra := map[string]interface{}{}
	extra["Image"] = clusterv1beta1.OdooImageSpec{
		Registry:   instance.Annotations["cluster.odoo.io/registry"],
		Repository: instance.Annotations["cluster.odoo.io/repository"],
		Trackname:  instance.Labels["app.kubernetes.io/track"],
		Version:    instance.Labels["app.kubernetes.io/version"],
		Secret:     instance.Annotations["cluster.odoo.io/pull-secret"],
	}
	res, _, err := ctx.CreateOrUpdate(comp.templatePath, extra, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*appsv1.Deployment)
		existing := existingObj.(*appsv1.Deployment)
		// Copy the configuration Data over.
		existing.Spec = goal.Spec
		return nil
	})
	return res, err
}
