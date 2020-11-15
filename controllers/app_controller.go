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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	runnerv1alpha1 "github.com/lqshow/kubernetes-crd/api/v1alpha1"
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
	instanceLog := r.Log.WithValues("app", req.NamespacedName)

	// your logic here
	instance := &runnerv1alpha1.App{}
	instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, "app.runner.basebit.me")

	// Get app instance
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		instanceLog.Error(err, "unable to fetch app")
	} else {
		instanceLog.Info(instance.Status.Status)
	}

	// updating the status
	if instance.Status.Status == "" {
		instance.Status.Status = "Running"
		err := r.Status().Update(ctx, instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runnerv1alpha1.App{}).
		Complete(r)
}
