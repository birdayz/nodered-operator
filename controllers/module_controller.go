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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	noderednerdendev1 "github.com/birdayz/nodered-operator/api/v1"
	"github.com/birdayz/nodered-operator/pkg/nodered"
)

// ModuleReconciler reconciles a Module object
type ModuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nodered.nerden.de.github.com,resources=modules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nodered.nerden.de.github.com,resources=modules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nodered.nerden.de.github.com,resources=modules/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Module object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var instance noderednerdendev1.Module
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	nodeRedList := noderednerdendev1.NoderedList{}
	if err := r.List(ctx, &nodeRedList, client.InNamespace(req.Namespace), client.MatchingLabels(instance.Spec.Selector)); err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Found", "items", len(nodeRedList.Items))

	for _, item := range nodeRedList.Items {
		l.Info("Fetch secret for", "instance", item.Name)

		var secret corev1.Secret
		secretName := types.NamespacedName{
			Name:      item.Name + "-nodered" + "-operator",
			Namespace: req.Namespace,
		}

		err := r.Get(ctx, secretName, &secret)
		if err != nil {
			if errors.IsNotFound(err) {
				return ctrl.Result{}, fmt.Errorf("could not find operator secret for instance: %w", err)
			} else {
				return ctrl.Result{}, err
			}
		}

		username := string(secret.Data["username"])
		password := string(secret.Data["password"])
		l.Info("Got secret", "username", username, "password", password)

		client := nodered.NewClient(fmt.Sprintf("%s-nodered.%s.svc.cluster.local", item.Name, item.Namespace), 1881, username, password)

		_, err = client.CreateModule(instance.Spec.PackageName)
		if err != nil {
			if nrError, ok := err.(*nodered.NodeRedError); ok {
				if nrError.Code == nodered.ErrorCodeModuleAlreadyLoaded {
					l.Info("Module already loaded, doing nothing")
				}
			} else {
				return ctrl.Result{}, err
			}
		}

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&noderednerdendev1.Module{}).
		Complete(r)
}
