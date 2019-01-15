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
	"bytes"
	ini "gopkg.in/ini.v1"

	"github.com/blaggacao/ridecell-operator/pkg/components"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

type configmapComponent struct {
	templatePath string
}

func NewConfigMap(templatePath string) *configmapComponent {
	return &configmapComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=configmaps/status,verbs=get;update;patch
func (_ *configmapComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&corev1.ConfigMap{},
	}
}

func (_ *configmapComponent) IsReconcilable(_ *components.ComponentContext) bool {
	// ConfigMaps have no dependencies, always reconcile.
	return true
}

func (comp *configmapComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	// instance := ctx.Top.(*clusterv1beta1.OdooVersion)
	// trackinstance := nil
	// clusterinstance := nil

	mergedConfig := map[string]clusterv1beta1.ConfigValue{}

	// Create config.ini
	cfg := ini.Empty()
	marshallConfig(cfg, mergedConfig, "")
	var buf bytes.Buffer
	cfg.WriteTo(&buf)

	// Set up the extra data map for the template.
	extra := map[string]interface{}{}
	extra["ConfigFile"] = buf.String()

	res, _, err := ctx.CreateOrUpdate(comp.templatePath, extra, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*corev1.ConfigMap)
		existing := existingObj.(*corev1.ConfigMap)
		// Copy the configuration Data over.
		existing.Data = goal.Data
		return nil
	})
	return res, err
}

func marshallConfig(cfg *ini.File, config map[string]clusterv1beta1.ConfigValue, section string) {
	for key, cfgValue := range config {
		if cfgValue.Section != nil {
			marshallConfig(cfg, cfgValue.Section, section+":"+key)
		} else {
			cfg.Section(section).NewKey(key, cfgValue.ToString())
		}
	}

}
