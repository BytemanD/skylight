package openstack

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DeleteOptions struct {
	options metav1.DeleteOptions
}

func (o *DeleteOptions) Options() metav1.DeleteOptions {
	return o.options
}
func (o *DeleteOptions) SetDryRun() {
	o.options.DryRun = append(o.options.DryRun, metav1.DryRunAll)
}
func (o *DeleteOptions) SetGracePeriodSeconds(g *int64) {
	o.options.GracePeriodSeconds = g
}

func NewDeleteOption() DeleteOptions {
	return DeleteOptions{}
}
