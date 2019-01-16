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
	"fmt"

	"github.com/blaggacao/ridecell-operator/pkg/components"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type metaComponent struct{}

func NewMeta() *metaComponent { return &metaComponent{} }

func (*metaComponent) WatchTypes() []runtime.Object { return []runtime.Object{} }

func (*metaComponent) IsReconcilable(_ *components.ComponentContext) bool { return true }

func (*metaComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {

	instance := ctx.Top.(*clusterv1beta1.OdooVersion)
	res, _, err := ctx.UpdateTopMeta(func(goalMeta *metav1.ObjectMeta) error {
		goalMeta.Labels = map[string]string{
			"cluster.odoo.io/name":         instance.Spec.Cluster,
			"cluster.odoo.io/track":        fmt.Sprintf("%v", instance.Spec.Track),
			"app.kubernetes.io/name":       "odooversion",
			"app.kubernetes.io/instance":   instance.Name,
			"app.kubernetes.io/component":  "cluster",
			"app.kubernetes.io/managed-by": "odoo-operator",
			"app.kubernetes.io/version":    instance.Spec.Version,
		}
		return nil
	})
	return res, err
}
