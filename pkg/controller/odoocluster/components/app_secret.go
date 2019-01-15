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
	"io/ioutil"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/blaggacao/ridecell-operator/pkg/components"
	clusterv1beta1 "github.com/xoe-labs/odoo-operator/pkg/apis/cluster/v1beta1"
)

const inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

type appSecretComponent struct {
	templatePath string
}

func NewAppSecret(templatePath string) *appSecretComponent {
	return &appSecretComponent{templatePath: templatePath}
}

// +kubebuilder:rbac:groups=,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=secrets/status,verbs=get;update;patch
func (_ *appSecretComponent) WatchTypes() []runtime.Object {
	return []runtime.Object{
		&corev1.Secret{},
	}
}

func (_ *appSecretComponent) IsReconcilable(_ *components.ComponentContext) bool {
	// ConfigMaps have no dependencies, always reconcile.
	return true
}

func (comp *appSecretComponent) Reconcile(ctx *components.ComponentContext) (reconcile.Result, error) {
	instance := ctx.Top.(*clusterv1beta1.OdooCluster)
	var upstream *corev1.Secret

	// Loan secret data from operator namespace
	if instance.Spec.AppSecret != "" {
		operatorNamespace := os.Getenv("NAMESPACE")
		if operatorNamespace == "" {
			var err error
			operatorNamespace, err = getInClusterNamespace()
			if err != nil {
				instance.SetStatusConditionOperatorNamespaceErrored()
				ctx.Logger.Error(err, "no namespace found")
				return reconcile.Result{}, err
			}
		}
		upstream = &corev1.Secret{}
		err := ctx.Get(ctx.Context, types.NamespacedName{Name: instance.Spec.AppSecret, Namespace: operatorNamespace}, upstream)
		if err != nil {
			if errors.IsNotFound(err) {
				instance.SetStatusConditionSecretLoaningNotFoundErrored()
			}
			return reconcile.Result{Requeue: true}, err
		}
	}

	res, _, err := ctx.CreateOrUpdate(comp.templatePath, nil, func(goalObj, existingObj runtime.Object) error {
		goal := goalObj.(*corev1.Secret)
		existing := existingObj.(*corev1.Secret)
		// Copy the secret Data over.
		existing.Data = goal.Data
		// Loan Odoo's adminpasswd from upstream, if available.
		if upstream != nil {
			_, ok := upstream.Data["adminpasswd"]
			if !ok {
				instance.SetStatusConditionSecretLoaningAdminPasswdNotFoundErrored()
				err := fmt.Errorf("app secret loaning failed: expected key not found")
				ctx.Logger.Error(err, "app secret loaning failed")
				return err
			}
			existing.Data["adminpasswd"] = upstream.Data["adminpasswd"]
			instance.SetStatusConditionSecretLoaningSuccessAppSecretLoaned()
		}
		existing.Type = goal.Type
		return nil
	})
	return res, err
}

func getInClusterNamespace() (string, error) {
	// Check whether the namespace file exists.
	// If not, we are not running in cluster so can't guess the namespace.
	_, err := os.Stat(inClusterNamespacePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("not running in-cluster, please specify $NAMESPACE")
	} else if err != nil {
		return "", fmt.Errorf("error checking namespace file: %v", err)
	}

	// Load the namespace file and return itss content
	namespace, err := ioutil.ReadFile(inClusterNamespacePath)
	if err != nil {
		return "", fmt.Errorf("error reading namespace file: %v", err)
	}
	return string(namespace), nil
}
