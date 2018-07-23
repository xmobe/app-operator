package cluster

import (
	"fmt"
	"path/filepath"

	// "github.com/sirupsen/logrus"
	api "github.com/xoe-labs/odoo-operator/pkg/apis/odoo/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func odooContainer(cr *api.OdooCluster, trackSpec *api.TrackSpec, tierSpec *api.TierSpec) v1.Container {

	args := getContainerArgs(tierSpec)
	ports := getContainerPorts(tierSpec)
	volumes := []v1.VolumeMount{
		{
			Name:      getVolumeName(cr, configVolName),
			MountPath: filepath.Dir(appConfigsPath),
			ReadOnly:  true,
		},
		{
			Name:      getVolumeName(cr, secretVolName),
			MountPath: filepath.Dir(appSecretsPath),
			ReadOnly:  true,
		},
	}

	for _, s := range cr.Spec.Volumes {
		volumes = append(volumes, v1.VolumeMount{
			Name:      getVolumeName(cr, s.Name),
			MountPath: filepath.Dir(getMountPath(s.Name)),
		})
	}

	c := v1.Container{
		Name:  getFullName(cr, trackSpec, tierSpec),
		Image: getImageName(&trackSpec.Image),
		Args:  args,
		Env: []v1.EnvVar{
			{
				Name:  envPGHOST,
				Value: cr.Spec.PgSpec.PgCluster.Host,
			},
			{
				Name:  envPGPASSFILE,
				Value: getSecretFile(appPsqlSecretKey),
			},
			{
				Name:  envODOORC,
				Value: appConfigsPath,
			},
			{
				Name:  envODOOPASSFILE,
				Value: getSecretFile(appAdminSecretKey),
			},
		},
		VolumeMounts: volumes,
		Ports:        ports,
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: v1.TerminationMessageReadFile,
		ImagePullPolicy:          v1.PullAlways,
	}
	switch tierSpec.Name {
	case api.ServerTier:
		c.LivenessProbe = &v1.Probe{
			Handler: v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{
						"curl",
						"--connect-timeout", "5",
						"--max-time", "10",
						"-k", "-s",
						fmt.Sprintf("https://localhost:%d/web", clientPort),
					},
				},
			},
			InitialDelaySeconds: 10,
			TimeoutSeconds:      10,
			PeriodSeconds:       60,
			FailureThreshold:    3,
			SuccessThreshold:    1,
		}
		c.ReadinessProbe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path:   "/web",
					Port:   intstr.FromInt(clientPort),
					Scheme: v1.URISchemeHTTPS,
				},
			},
			InitialDelaySeconds: 10,
			TimeoutSeconds:      10,
			PeriodSeconds:       10,
			FailureThreshold:    3,
			SuccessThreshold:    1,
		}
	case api.LongpollingTier:
		c.LivenessProbe = &v1.Probe{
			Handler: v1.Handler{
				Exec: &v1.ExecAction{
					Command: []string{
						"curl",
						"--connect-timeout", "5",
						"--max-time", "10",
						"-k", "-s",
						fmt.Sprintf("https://localhost:%d/web", longpollingPort),
					},
				},
			},
			InitialDelaySeconds: 10,
			TimeoutSeconds:      10,
			PeriodSeconds:       60,
			FailureThreshold:    3,
			SuccessThreshold:    1,
		}
		c.ReadinessProbe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path:   "/web",
					Port:   intstr.FromInt(longpollingPort),
					Scheme: v1.URISchemeHTTPS,
				},
			},
			InitialDelaySeconds: 10,
			TimeoutSeconds:      10,
			PeriodSeconds:       10,
			FailureThreshold:    3,
			SuccessThreshold:    1,
		}
	}
	return c
}

func getContainerArgs(tierSpec *api.TierSpec) []string {
	switch tierSpec.Name {
	case api.ServerTier:
		return []string{"--config", appConfigsPath}
		// return []string{"--config", appConfigsPath, "--tier", api.ServerTier}
	case api.CronTier:
		return []string{"--config", appConfigsPath}
		// return []string{"--config", appConfigsPath, "--tier", api.CronTier}
	case api.BackgroundTier:
		return []string{"--config", appConfigsPath}
		// return []string{"--config", appConfigsPath, "--tier", api.BackgroundTier}
	case api.LongpollingTier:
		return []string{"--config", appConfigsPath}
		// return []string{"--config", appConfigsPath, "--tier", api.LongpollingTier}
	}
	return nil
}

func getContainerPorts(tierSpec *api.TierSpec) []v1.ContainerPort {
	switch tierSpec.Name {
	case api.ServerTier:
		return []v1.ContainerPort{{
			Name:          clientPortName,
			ContainerPort: int32(clientPort),
			Protocol:      v1.ProtocolTCP,
		}}
	case api.CronTier:
		return []v1.ContainerPort{}
	case api.BackgroundTier:
		return []v1.ContainerPort{}
	case api.LongpollingTier:
		return []v1.ContainerPort{{
			Name:          longpollingPortName,
			ContainerPort: int32(longpollingPort),
			Protocol:      v1.ProtocolTCP,
		}}
	}
	return nil
}
