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
package odooversion

import (
	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/blaggacao/ridecell-operator/pkg/components"
	odooversioncomponents "github.com/xoe-labs/odoo-operator/pkg/controller/odooversion/components"
)

// Add creates a new OdooVersion Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
//
// +kubebuilder:rbac:groups=cluster.odoo.io,resources=odooversions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cluster.odoo.io,resources=odooversions/status,verbs=get;update;patch
func Add(mgr manager.Manager) error {
	_, err := components.NewReconciler("odooversion.controller", mgr, &clusterv1beta1.OdooVersion{}, Templates, []components.Component{
		// Set default values.
		// odooversioncomponents.NewDefaults(),

		// Set Top-level components
		odooversioncomponents.NewConfigMap("configmap.yml.tpl"),

		// Web Server components (app.kubernetes.io/component = web)
		odooversioncomponents.NewDeployment("web/deployment.yml.tpl"),
		// odooversioncomponents.NewService("web/service.yml.tpl"),

		// Longpolling components (app.kubernetes.io/component = longpolling)
		odooversioncomponents.NewDeployment("longpolling/deployment.yml.tpl"),
		// odooversioncomponents.NewService("longpolling/service.yml.tpl"),

		// Cron components (app.kubernetes.io/component = cron)
		odooversioncomponents.NewDeployment("cron/deployment.yml.tpl"),

		// Remover acting upon finalizers of deleted instances
		// odooversioncomponents.NewRemover(),

		// Done on the cluster controller as owner of the ingress resource
		// // L7 instance routing
		// // Keep it at the end of this block to consume a consistent final state
		// odooversioncomponents.NewRouter("sync-migrator.yml.tpl"),
	})
	return err
}
