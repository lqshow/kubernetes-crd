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
	"github.com/lqshow/kubernetes-crd/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runnerv1alpha1 "github.com/lqshow/kubernetes-crd/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// FuwuReconciler reconciles a Fuwu object
type FuwuReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=runner.basebit.me,resources=fuwus,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=runner.basebit.me,resources=fuwus/status,verbs=get;update;patch

func (r *FuwuReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	instanceLog := r.Log.WithValues("fuwu", req.NamespacedName)

	// your logic here
	instance := &runnerv1alpha1.Fuwu{}

	// Get fuwu instance
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		instanceLog.Error(err, "unable to fetch fuwu")
	} else {
		instanceLog.Info(instance.Spec.Name, instance.Spec.Description, instance.Status)
	}

	// updating fuwu status
	if instance.Status.Status == "" {
		instance.Status.Status = "Running"
		if err := r.Status().Update(ctx, instance); err != nil {
			instanceLog.Error(err, "unable to update fuwu status")
		}
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		finalizerName := "fuwu.runner.basebit.me"
		if !util.InStringArray(instance.ObjectMeta.Finalizers, finalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				instanceLog.Error(err, "unable to update fuwu")
				return ctrl.Result{}, err
			}
		}
	} else {
		// TODO
	}

	// Delete fuwu
	//time.Sleep(time.Second * 10)
	//if err := r.Delete(ctx, instance); err != nil {
	//	log.Error(err, "unable to delete fuwu", "fuwu", instance)
	//}

	return ctrl.Result{}, nil
}

func (r *FuwuReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runnerv1alpha1.Fuwu{}).
		Complete(r)
}
