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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/blaggacao/ridecell-operator/pkg/components"
)

// OdooCluster Status ...
func (c *OdooCluster) GetStatus() components.Status {
	return c.Status
}

func (c *OdooCluster) SetStatus(status components.Status) {
	c.Status = status.(OdooClusterStatus)
}

func (c *OdooCluster) SetErrorStatus(errorMsg string) {
	// c.Status.Status = StatusError
	// c.Status.Message = errorMsg
}

func (_ *OdooCluster) NewStatusCondition(
	condType OdooClusterStatusConditionType, status corev1.ConditionStatus,
	reason, message string) *OdooClusterStatusCondition {
	return &OdooClusterStatusCondition{condType, StatusCondition{
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}}
}

func (c *OdooCluster) GetStatusCondition(
	condType OdooClusterStatusConditionType) *OdooClusterStatusCondition {
	for i := range c.Status.Conditions {
		cond := c.Status.Conditions[i]
		if cond.Type == condType {
			return &cond
		}
	}
	return nil
}

func (c *OdooCluster) SetStatusCondition(condition OdooClusterStatusCondition) {
	currentCond := c.GetStatusCondition(condition.Type)
	if currentCond != nil && currentCond.Status == condition.Status && currentCond.Reason == condition.Reason {
		return
	}
	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := c.filterOutCondition(condition.Type)
	c.Status.Conditions = append(newConditions, condition)
}

func (c *OdooCluster) RemoveStatusCondition(condType OdooClusterStatusConditionType) {
	c.Status.Conditions = c.filterOutCondition(condType)
}

func (c *OdooCluster) filterOutCondition(condType OdooClusterStatusConditionType) []OdooClusterStatusCondition {
	var newConditions []OdooClusterStatusCondition
	for _, cond := range c.Status.Conditions {
		if cond.Type == condType {
			continue
		}
		newConditions = append(newConditions, cond)
	}
	return newConditions
}

// OdooTrack Status ...
func (t *OdooTrack) GetStatus() components.Status {
	return t.Status
}

func (t *OdooTrack) SetStatus(status components.Status) {
	t.Status = status.(OdooTrackStatus)
}
func (s *OdooTrack) SetErrorStatus(errorMsg string) {
	// t.Status.Status = StatusError
	// t.Status.Message = errorMsg
}

// OdooVersion Status ...
func (v *OdooVersion) GetStatus() components.Status {
	return v.Status
}

func (v *OdooVersion) SetStatus(status components.Status) {
	v.Status = status.(OdooVersionStatus)
}
func (v *OdooVersion) SetErrorStatus(errorMsg string) {
	// v.Status.Status = StatusError
	// v.Status.Message = errorMsg
}

func (_ *OdooVersion) NewStatusCondition(
	condType OdooVersionStatusConditionType, status corev1.ConditionStatus,
	reason, message string) *OdooVersionStatusCondition {
	return &OdooVersionStatusCondition{condType, StatusCondition{
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}}
}

func (v *OdooVersion) GetStatusCondition(
	condType OdooVersionStatusConditionType) *OdooVersionStatusCondition {
	for i := range v.Status.Conditions {
		cond := v.Status.Conditions[i]
		if cond.Type == condType {
			return &cond
		}
	}
	return nil
}

func (v *OdooVersion) SetStatusCondition(condition OdooVersionStatusCondition) {
	currentCond := v.GetStatusCondition(condition.Type)
	if currentCond != nil && currentCond.Status == condition.Status && currentCond.Reason == condition.Reason {
		return
	}
	// Do not update lastTransitionTime if the status of the condition doesn't change.
	if currentCond != nil && currentCond.Status == condition.Status {
		condition.LastTransitionTime = currentCond.LastTransitionTime
	}
	newConditions := v.filterOutCondition(condition.Type)
	v.Status.Conditions = append(newConditions, condition)
}

func (v *OdooVersion) RemoveStatusCondition(condType OdooVersionStatusConditionType) {
	v.Status.Conditions = v.filterOutCondition(condType)
}

func (v *OdooVersion) filterOutCondition(condType OdooVersionStatusConditionType) []OdooVersionStatusCondition {
	var newConditions []OdooVersionStatusCondition
	for _, cond := range v.Status.Conditions {
		if cond.Type == condType {
			continue
		}
		newConditions = append(newConditions, cond)
	}
	return newConditions
}
