package cluster

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/klog"
	providerv1 "sigs.k8s.io/cluster-api-provider-openstack/pkg/apis/openstackproviderconfig/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/cloud/openstack"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	client "sigs.k8s.io/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
)

type Actuator struct {
	params         openstack.ActuatorParams
	clustersGetter client.ClustersGetter
}

// NewActuator creates a new Actuator
func NewActuator(params openstack.ActuatorParams) (*Actuator, error) {
	res := &Actuator{params: params}
	return res, nil
}

func (a *Actuator) Reconcile(cluster *clusterv1.Cluster) error {
	klog.Infof("Reconciling cluster %v.", cluster.Name)

	// Load provider config.
	_, err := providerv1.ClusterSpecFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		return errors.Errorf("failed to load cluster provider spec: %v", err)
	}

	// Load provider status.
	_, err = providerv1.ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	/* Uncomment when the clusterGetter is back to working
	defer func() {
		if err := a.storeClusterStatus(cluster, status); err != nil {
			klog.Errorf("failed to store provider status for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
		}
	}()*/
	return nil
}

// Delete deletes a cluster and is invoked by the Cluster Controller
func (a *Actuator) Delete(cluster *clusterv1.Cluster) error {
	klog.Infof("Deleting cluster %v.", cluster.Name)

	// Load provider config.
	_, err := providerv1.ClusterSpecFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		return errors.Errorf("failed to load cluster provider config: %v", err)
	}

	// Load provider status.
	_, err = providerv1.ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
	if err != nil {
		return errors.Errorf("failed to load cluster provider status: %v", err)
	}

	// Delete other things

	return nil
}

func (a *Actuator) storeClusterStatus(cluster *clusterv1.Cluster, status *providerv1.OpenstackClusterProviderStatus) error {
	clusterClient := a.clustersGetter.Clusters(cluster.Namespace)

	ext, err := providerv1.EncodeClusterStatus(status)
	if err != nil {
		return fmt.Errorf("failed to update cluster status for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
	}

	cluster.Status.ProviderStatus = ext

	if _, err := clusterClient.UpdateStatus(cluster); err != nil {
		return fmt.Errorf("failed to update cluster status for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
	}

	return nil
}
