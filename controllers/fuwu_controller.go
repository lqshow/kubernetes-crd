/*


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
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runnerv1alpha1 "github.com/lqshow/kubernetes-crd/api/v1alpha1"
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
	_ = r.Log.WithValues("fuwu", req.NamespacedName)

	// your logic here
	fuwu := &runnerv1alpha1.Fuwu{}
	fuwu.ObjectMeta.Finalizers = append(fuwu.ObjectMeta.Finalizers, "fuwu.runner.basebit.me")

	// Get fuwu
	if err := r.Get(ctx, req.NamespacedName, fuwu); err != nil {
		fmt.Println(err, "unable to fetch fuwu")
	} else {
		fmt.Println(fuwu.Spec.Name, fuwu.Spec.Description)
	}

	// Update fuwu
	fuwu.Status.Status = "Running"
	if err := r.Status().Update(ctx, fuwu); err != nil {
		log.Error(err, "unable to update fuwu status")
	}

	// Delete fuwu
	//time.Sleep(time.Second * 10)
	//if err := r.Delete(ctx, fuwu); err != nil {
	//	log.Error(err, "unable to delete fuwu", "fuwu", fuwu)
	//}

	return ctrl.Result{}, nil
}

func (r *FuwuReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runnerv1alpha1.Fuwu{}).
		Complete(r)
}
