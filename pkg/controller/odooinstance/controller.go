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
package odooinstance

import (
	instancev1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/instance/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/blaggacao/ridecell-operator/pkg/components"
	odooinstancecomponents "github.com/xoe-labs/odoo-operator/pkg/controller/odooinstance/components"
)

// Add creates a new OdooInstance Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// +kubebuilder:rbac:groups=instance.odoo.io,resources=odooinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=instance.odoo.io,resources=odooinstances/status,verbs=get;update;patch
func Add(mgr manager.Manager) error {
	_, err := components.NewReconciler("odooinstance.controller", mgr, &instancev1beta1.OdooInstance{}, Templates, []components.Component{
		// Set default values and meta data.
		odooinstancecomponents.NewDefaults(),
		odooinstancecomponents.NewMeta(),

		// Top-level components
		odooinstancecomponents.NewInitializer("initializer.yml.tpl"),
		odooinstancecomponents.NewCopier("copier.yml.tpl"),

		// Secondary components
		// odooinstancecomponents.NewSyncMigrator("sync-migrator.yml.tpl"),
		// odooinstancecomponents.NewAsynchMigrator("asynch-migrator.yml.tpl"),

		// Backup job components
		odooinstancecomponents.NewBackuper("backuper.yml.tpl"),

		// Remover acting upon finalizers of deleted instances
		// odooinstancecomponents.NewRemover(),
	})
	return err
}
