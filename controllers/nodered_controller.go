/*
Copyright 2021.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	noderednerdendev1 "github.com/birdayz/nodered-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NoderedReconciler reconciles a Nodered object
type NoderedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nodered.nerden.de.github.com,resources=nodereds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nodered.nerden.de.github.com,resources=nodereds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nodered.nerden.de.github.com,resources=nodereds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Nodered object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NoderedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Test", "name", req.NamespacedName)

	resourceName := req.NamespacedName.Name + "-nodered"

	sts := appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: req.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"nodered.nerden.de/instance": req.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"nodered.nerden.de/instance": req.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "nodered",
							Image:           "nodered/node-red",
							ImagePullPolicy: corev1.PullIfNotPresent,
							// Env: []corev1.EnvVar{
							// 	{
							// 		Name:  "APISERVER_URL",
							// 		Value: "https://" + instance.Spec.Apiserver.Restful.Host,
							// 	},
							// },
						},
					},
				},
			},
		},
	}

	if err := r.Get(ctx, req.NamespacedName, &sts); err != nil {
		if errors.IsNotFound(err) {
			err = r.Create(ctx, &sts)
			if err != nil {
				return ctrl.Result{}, err
			}
			l.Info("Created StatefulSet", "name", req.Name, "namespace", req.Namespace)
		} else {
			return ctrl.Result{}, err
		}
	}

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NoderedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&noderednerdendev1.Nodered{}).
		Complete(r)
}
