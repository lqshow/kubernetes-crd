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
	"fmt"

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
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get

func (r *FuwuReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("fuwu", req.NamespacedName)

	// your logic here
	instance := &runnerv1alpha1.Fuwu{}

	// Get fuwu instance
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		log.Error(err, "unable to fetch fuwu")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.V(1).Info("Getting fuwu instance", "spec", instance.Spec, "status", instance.Status)

	// updating fuwu status
	if instance.Status.Status == "" {
		instance.Status.Status = "Running"
		if err := r.Status().Update(ctx, instance); err != nil {
			log.Error(err, "unable to update fuwu status")
		}
	}

	finalizerName := "fuwu.runner.basebit.me"
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !util.InStringArray(instance.ObjectMeta.Finalizers, finalizerName) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				log.Error(err, "unable to update fuwu")
				return ctrl.Result{}, err
			}
		}
	} else {
		if util.InStringArray(instance.ObjectMeta.Finalizers, finalizerName) {
			// TODO: delete additional resources
			if err := r.deleteAppResources(req, ctx, instance); err != nil {
				return ctrl.Result{}, err
			}

			instance.ObjectMeta.Finalizers = util.RemoveString(instance.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				log.Error(err, "unable to update fuwu")
				return ctrl.Result{}, err
			}
		}
	}

	if err := r.reconcileApp(req, ctx, instance); err != nil {
		log.Error(err, "Creating App error")
		return ctrl.Result{}, err
	}

	//var pods corev1.PodList
	//if err := r.List(ctx, &pods, client.InNamespace(req.Namespace)); err != nil {
	//	log.Error(err, "unable to fetch Pod")
	//	return ctrl.Result{}, client.IgnoreNotFound(err)
	//}
	//log.V(1).Info("Getting pod", "pod spec", pods)

	// Delete fuwu
	//time.Sleep(time.Second * 10)
	//if err := r.Delete(ctx, instance); err != nil {
	//	log.Error(err, "unable to delete fuwu", "fuwu", instance)
	//}

	return ctrl.Result{}, nil
}

func (r *FuwuReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Watch for changes to Fuwu
		For(&runnerv1alpha1.Fuwu{}).
		// Watch App created by Fuwu
		Owns(&runnerv1alpha1.App{}).
		Complete(r)
}

func (r *FuwuReconciler) deleteAppResources(req ctrl.Request, ctx context.Context, instance *runnerv1alpha1.Fuwu) error {
	// 删除 Fuwu 时，必须将 Fuwu 资源的 Selector 指向的所有 App 都删除
	// 需检查 Spec.Selector 是否不为空
	if instance.Spec.Selector == "" {
		return nil
	}
	// Spec.Selector 如果不为空

	var apps runnerv1alpha1.AppList
	if err := r.List(ctx, &apps, client.InNamespace(req.Namespace)); err != nil {

	}

	for _, app := range apps.Items {
		if err := r.Delete(ctx, &app); err != nil {
			return fmt.Errorf("error occurred delete app: namespace=%s, name=%s", app.Namespace, app.Name)
		}
	}

	return nil
}

func (r *FuwuReconciler) reconcileApp(req ctrl.Request, ctx context.Context, instance *runnerv1alpha1.Fuwu) error {
	//for _, app := range instance.Spec.Apps {
	//	_ = app
	//}
	return nil
}
