package config

import (
	"testing"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
)

func TestFindDefaultManagedProxyConfigurationName(t *testing.T) {
	cases := []struct {
		name               string
		cma                *addonv1alpha1.ClusterManagementAddOn
		expectedConfigName string
	}{
		{
			name: "no config",
			cma:  &addonv1alpha1.ClusterManagementAddOn{},
		},
		{
			name: "non proxy.open-cluster-management.io",
			cma: &addonv1alpha1.ClusterManagementAddOn{
				Spec: addonv1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []addonv1alpha1.ConfigMeta{
						{
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
								Group:    "test.io",
								Resource: "tests",
							},
						},
					},
				},
			},
		},
		{
			name: "non managed proxy config",
			cma: &addonv1alpha1.ClusterManagementAddOn{
				Spec: addonv1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []addonv1alpha1.ConfigMeta{
						{
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
								Group:    "proxy.open-cluster-management.io",
								Resource: "tests",
							},
						},
					},
				},
			},
		},
		{
			name: "no defautl config",
			cma: &addonv1alpha1.ClusterManagementAddOn{
				Spec: addonv1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []addonv1alpha1.ConfigMeta{
						{
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
								Group:    "proxy.open-cluster-management.io",
								Resource: "managedproxyconfigurations",
							},
						},
					},
				},
			},
		},
		{
			name: "has managed proxy config",
			cma: &addonv1alpha1.ClusterManagementAddOn{
				Spec: addonv1alpha1.ClusterManagementAddOnSpec{
					SupportedConfigs: []addonv1alpha1.ConfigMeta{
						{
							ConfigGroupResource: addonv1alpha1.ConfigGroupResource{
								Group:    "proxy.open-cluster-management.io",
								Resource: "managedproxyconfigurations",
							},
							DefaultConfig: &addonv1alpha1.ConfigReferent{
								Name: "cluster-proxy",
							},
						},
					},
				},
			},
			expectedConfigName: "cluster-proxy",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := FindDefaultManagedProxyConfigurationName(c.cma)
			if actual != c.expectedConfigName {
				t.Errorf("expected %q, but %q", c.expectedConfigName, actual)
			}
		})
	}

}
