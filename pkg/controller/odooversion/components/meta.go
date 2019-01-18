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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type metaComponent struct{}

func NewMeta() *metaComponent { return &metaComponent{} }

func (*metaComponent) WatchTypes() []runtime.Object { return []runtime.Object{} }

func (*metaComponent) IsReconcilable(_ *components.ComponentContext) bool { return true }

func (*metaComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	instance := ctx.Top.(*clusterv1beta1.OdooVersion)

	// Fetch OdooCluster
	clusterinstance := &clusterv1beta1.OdooCluster{}
	err := ctx.Get(ctx.Context, client.ObjectKey{
		Name:      instance.Spec.Cluster,
		Namespace: instance.GetNamespace(),
	}, clusterinstance)
	if err != nil {
		return reconcile.Result{}, err
	}
	// Fetch OdooTrack
	listObjTrack := &clusterv1beta1.OdooTrackList{}
	res, obj, err := ctx.GetOne(listObjTrack, map[string]string{
		"cluster.odoo.io/part-of-cluster": instance.Spec.Cluster,
		"app.kubernetes.io/track":         fmt.Sprintf("%v", instance.Spec.Track),
	})
	if err != nil || obj == nil {
		return res, err
	}
	odooTrack := obj.(*clusterv1beta1.OdooTrack)
	res, _, err = ctx.UpdateTopMeta(func(goalMeta *metav1.ObjectMeta) error {
		goalMeta.Labels = map[string]string{
			"cluster.odoo.io/part-of-cluster": instance.Spec.Cluster,
			"cluster.odoo.io/part-of-track":   odooTrack.Name,
			"app.kubernetes.io/name":          "version",
			"app.kubernetes.io/instance":      instance.Name,
			"app.kubernetes.io/component":     "operator",
			"app.kubernetes.io/managed-by":    "odoo-operator",
			"app.kubernetes.io/version":       instance.Spec.Version,
			"app.kubernetes.io/part-of":       odooTrack.Name,
			"app.kubernetes.io/track":         fmt.Sprintf("%v", instance.Spec.Track),
		}
		goalMeta.Annotations = map[string]string{
			"cluster.odoo.io/registry":   clusterinstance.Spec.Image.Registry,
			"cluster.odoo.io/pull-seret": clusterinstance.Spec.Image.Secret,
			"cluster.odoo.io/repository": clusterinstance.Spec.Image.Repo,
		}
		return nil
	})
	return res, err
}
