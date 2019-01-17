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
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type odooTrackComponent struct {
	templatePath string
}

func NewOdooTrack(templatePath string) *odooTrackComponent {
	return &odooTrackComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=cluster.odoo.io,resources=odootracks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.odoo.io,resources=odootracks/status,verbs=get;update;patch
func (_ *odooTrackComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&clusterv1beta1.OdooTrack{},
	}
}

func (_ *odooTrackComponent) IsReconcilable(_ *components.ComponentContext) bool {
	return true
}

func (comp *odooTrackComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	instance := ctx.Top.(*clusterv1beta1.OdooCluster)
	var res reconcile.Result
	var err error
	// Iterate over all tracks and create CR
	for _, track := range instance.Spec.Tracks {
		// Prepare extra data with OdooInstanceList
		extra := map[string]interface{}{}
		extra["Name"] = fmt.Sprintf("%s-%v", instance.Name, track.Track)
		extra["Track"] = fmt.Sprintf("%v", track.Track)
		res, _, err = ctx.CreateOrUpdate(comp.templatePath, extra, func(goalObj, existingObj runtime.Object) error {
			goal := goalObj.(*clusterv1beta1.OdooTrack)
			existing := existingObj.(*clusterv1beta1.OdooTrack)
			// Copy over the (optional) config (without yaml indirection)
			goal.Spec.Config = track.Config
			// Copy the Spec over.
			existing.Spec = goal.Spec
			return nil
		})
		if err != nil {
			return res, err
		}
	}
	return res, err

}
