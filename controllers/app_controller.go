/*
Copyright 2020 LQ.

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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	runnerv1alpha1 "github.com/lqshow/kubernetes-crd/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=runner.basebit.me,resources=apps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=runner.basebit.me,resources=apps/status,verbs=get;update;patch

func (r *AppReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("app", req.NamespacedName)

	// Fetch the App instance
	instance := &runnerv1alpha1.App{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.V(1).Info("Getting app instance", "spec", instance.Spec, "status", instance.Status, "request name", req.Name)

	if !instance.DeletionTimestamp.IsZero() {
		log.V(1).Info("Getting deleted app instance")
		return ctrl.Result{}, nil
	}

	// If no phase set, default to pending (the initial phase):
	if instance.Status.Phase == "" {
		instance.Status.Phase = runnerv1alpha1.PhasePending
	}

	// Sync app instance status
	if err := r.syncAppStatus(ctx, instance); err != nil {
		log.Error(err, "Sync app status error")
		return ctrl.Result{}, err
	}

	// Reconcile app instance
	if err := r.reconcileInstance(ctx, instance); err != nil {
		log.Error(err, "Reconcile app instance error")
		return ctrl.Result{}, err
	}

	//instance.OwnerReferences =
	//instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, "app.runner.basebit.me")
	return ctrl.Result{}, nil
}

func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Watch for changes to App
		For(&runnerv1alpha1.App{}).
		// Watch Deployment created by App
		Owns(&appsv1.Deployment{}).
		// Watch Service created by App
		Owns(&corev1.Service{}).
		Complete(r)
}

func (r *AppReconciler) syncAppStatus(ctx context.Context, instance *runnerv1alpha1.App) error {
	r.Log.V(1).Info("Geting app instance status", "Phase", runnerv1alpha1.PhasePending)

	switch instance.Status.Phase {
	case runnerv1alpha1.PhasePending:
		// TODO: create service
		// TODO: create deployment
		instance.Status.Phase = runnerv1alpha1.PhaseRunning
	case runnerv1alpha1.PhaseRunning:
		// TODO: check deployment status
	case runnerv1alpha1.PhaseDone:
		return nil
	default:
		return nil
	}

	err := r.Status().Update(ctx, instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *AppReconciler) reconcileInstance(ctx context.Context, instance *runnerv1alpha1.App) error {

	//instance.Spec.AppNode
	deploy, err := buildDeployment(instance)
	if err != nil {
		return err
	}

	// Set At instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deploy, r.Scheme); err != nil {
		return err
	}

	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      deploy.Name,
		Namespace: deploy.Namespace,
	}, found)

	// Try to see if the deployment already exists and if not
	// (which we expect) then create a one-shot deployment as per spec
	if err != nil && errors.IsNotFound(err) {
		if err := r.Create(ctx, deploy); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		// Update the found object and write the result back if there are any changes
		found.Spec = deploy.Spec
		if err := r.Update(ctx, found); err != nil {
			return err
		}
	}

	return nil
}

func buildDeployment(instance *runnerv1alpha1.App) (*appsv1.Deployment, error) {
	labels := instance.Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["app.runner.basebit.me/app"] = instance.Name

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.GetNamespace(),
			Labels:    labels,
		},
		//Spec: *instance.Spec.AppNode.Template,
	}

	return deploy, nil
}
