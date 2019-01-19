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
)

// OdooInstanceStatusConditionType ...
type OdooInstanceStatusConditionType string

const (
	// OdooInstanceStatusConditionTypeCreated ...
	OdooInstanceStatusConditionTypeCreated OdooInstanceStatusConditionType = "Created"
	// OdooInstanceStatusConditionTypeMigrated ...
	OdooInstanceStatusConditionTypeMigrated OdooInstanceStatusConditionType = "Migrated"
	// OdooInstanceStatusConditionTypeBackuped ...
	OdooInstanceStatusConditionTypeBackuped OdooInstanceStatusConditionType = "Backuped"
	// OdooInstanceStatusConditionTypeReconciled ...
	OdooInstanceStatusConditionTypeReconciled OdooInstanceStatusConditionType = "Reconciled"
	// OdooInstanceStatusConditionTypeMaintaining ...
	OdooInstanceStatusConditionTypeMaintaining OdooInstanceStatusConditionType = "Maintaining"
)

func (i *OdooInstance) SetStatusConditionBackupCronJobSecheduleBackuped() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeBackuped,
		corev1.ConditionFalse,
		"BackupCronJobSechedule",
		"A backup CronJob schedule has been created to backup this database instance.")
	i.SetStatusCondition(*condition)
}

// copier:
func (i *OdooInstance) SetStatusConditionCopyJobCreationCreated() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeCreated,
		corev1.ConditionFalse,
		"CopyJobCreation",
		"A copier Job has been launched to copy and initialize this database instance.")
	i.SetStatusCondition(*condition)
}
func (i *OdooInstance) SetStatusConditionCopyJobSuccessCreated() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeCreated,
		corev1.ConditionTrue,
		"CopyJobCreation",
		"The database instance has been sucessfully created by a copier Job.")
	i.SetStatusCondition(*condition)
}

// initializer:
func (i *OdooInstance) SetStatusConditionInitJobCreationCreated() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeCreated,
		corev1.ConditionFalse,
		"InitJobCreation",
		"An initializer Job has been launched to initialize this database instance.")
	i.SetStatusCondition(*condition)
}
func (i *OdooInstance) SetStatusConditionInitJobSuccessCreated() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeCreated,
		corev1.ConditionTrue,
		"InitJobSuccess",
		"The database instance has been sucessfully created by an initializer Job.")
	i.SetStatusCondition(*condition)
}

// sync-migrator:
func (i *OdooInstance) SetStatusConditionSynchronousMigratorJobCreationMigrated() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeMigrated,
		corev1.ConditionFalse,
		"SynchronousMigratorJobCreation",
		"A synchronous migrator Job has been launched to migrate this database instance.")
	i.SetStatusCondition(*condition)
}
func (i *OdooInstance) SetStatusConditionSynchronousMigratorJobSuccessMigrated() {
	condition := i.NewStatusCondition(
		OdooInstanceStatusConditionTypeMigrated,
		corev1.ConditionTrue,
		"SynchronousMigratorJobSuccess",
		"The database instance has been sucessfully migrated by a synchronous migrator Job.")
	i.SetStatusCondition(*condition)
}
