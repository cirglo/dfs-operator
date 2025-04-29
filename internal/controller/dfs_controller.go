/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"github.com/cirglo/dfs-operator/assets"
	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	storagev1alpha1 "github.com/cirglo/dfs-operator/api/v1alpha1"
)

// DFSReconciler reconciles a DFS object
type DFSReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=storage.cirglo.com,resources=dfs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=storage.cirglo.com,resources=dfs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=storage.cirglo.com,resources=dfs/finalizers,verbs=update
// +kubebuilder:rbac:groups=storage.cirglo.com,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=storage.cirglo.com,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DFS object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *DFSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	cr := &storagev1alpha1.DFS{}
	err := r.Get(ctx, req.NamespacedName, cr)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "Operator resource not found")
			return ctrl.Result{}, err
		} else {
			logger.Error(err, "Failed to get Operator resource")
			return ctrl.Result{}, err
		}
	}

	nameNodeDeployment := &appsv1.Deployment{}
	err = r.Get(ctx, req.NamespacedName, nameNodeDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("NameNode Deployment not found, creating it")
			nameNodeDeployment := assets.GetNameNodeDeploymentFromFile()
			nameNodeDeployment.Namespace = req.Namespace
			nameNodeDeployment.Name = req.Name
			nameNodeDeployment.Spec.Replicas = cr.Spec.NumNameNodeServers
			nameNodeDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *cr.Spec.NameNodePort
			err = ctrl.SetControllerReference(cr, nameNodeDeployment, r.Scheme)
			if err != nil {
				logger.Error(err, "Failed to set controller reference for NameNode Deployment")
				return ctrl.Result{}, err
			}
			err = r.Create(ctx, nameNodeDeployment)
			if err != nil {
				logger.Error(err, "Failed to create NameNode Deployment")
				return ctrl.Result{}, err
			}

		} else {
			logger.Error(err, "Failed to get NameNode Deployment")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Updating NameNode Deployment")
		nameNodeDeployment.Spec.Replicas = cr.Spec.NumNameNodeServers
		nameNodeDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *cr.Spec.NameNodePort
		err = ctrl.SetControllerReference(cr, nameNodeDeployment, r.Scheme)
		if err != nil {
			logger.Error(err, "Failed to set controller reference for NameNode Deployment")
			return ctrl.Result{}, err
		}
		err = r.Update(ctx, nameNodeDeployment)
		if err != nil {
			logger.Error(err, "Failed to update NameNode Deployment")
			return ctrl.Result{}, err
		}
	}

	dataNodeStatefulSet := &appsv1.StatefulSet{}
	err = r.Get(ctx, req.NamespacedName, dataNodeStatefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("DataNode StatefulSet not found, creating it")
			dataNodeStatefulSet := assets.GetDataNodeStatefulSetFromFile()
			dataNodeStatefulSet.Namespace = req.Namespace
			dataNodeStatefulSet.Name = req.Name
			dataNodeStatefulSet.Spec.Replicas = cr.Spec.NumDataNodeServers
			dataNodeStatefulSet.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *cr.Spec.DataNodePort
			err = ctrl.SetControllerReference(cr, dataNodeStatefulSet, r.Scheme)
			if err != nil {
				logger.Error(err, "Failed to set controller reference for DataNode StatefulSet")
				return ctrl.Result{}, err
			}
			err = r.Create(ctx, dataNodeStatefulSet)
			if err != nil {
				logger.Error(err, "Failed to create DataNode StatefulSet")
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "Failed to get DataNode StatefulSet")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Updating DataNode StatefulSet")
		dataNodeStatefulSet.Spec.Replicas = cr.Spec.NumDataNodeServers
		dataNodeStatefulSet.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *cr.Spec.DataNodePort
		err = r.Update(ctx, dataNodeStatefulSet)
		if err != nil {
			logger.Error(err, "Failed to update DataNode StatefulSet")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DFSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storagev1alpha1.DFS{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
