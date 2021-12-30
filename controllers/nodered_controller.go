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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	noderednerdendev1 "github.com/birdayz/nodered-operator/api/v1"
	"github.com/sethvargo/go-password/password"
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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims/status,verbs=get

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

	compareOpts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
		patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus(),
	}

	instance := &noderednerdendev1.Nodered{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	name := types.NamespacedName{
		Name:      req.Name + "-nodered",
		Namespace: req.Namespace,
	}

	// Secret

	// Generate secret. The secret is truth for us, so we'll use it to set the
	// operator-account's password in the node-red resource.

	var operatorSecret corev1.Secret
	{
		secretName := types.NamespacedName{
			Name:      req.Name + "-nodered" + "-operator",
			Namespace: req.Namespace,
		}

		if err := r.Get(ctx, secretName, &operatorSecret); err != nil {
			if errors.IsNotFound(err) {

				//if it does not exist, generate a password and create it
				password, err := password.Generate(32, 10, 10, false, false)
				if err != nil {
					return ctrl.Result{}, fmt.Errorf("failed to generate random password: %w", err)
				}

				operatorSecret = corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      secretName.Name,
						Namespace: req.Namespace,
					},
					StringData: map[string]string{
						"username": "nodered-operator",
						"password": string(password),
					},
				}

				if err := controllerutil.SetControllerReference(instance, &operatorSecret, r.Scheme); err != nil {
					return ctrl.Result{}, err
				}

				err = r.Create(ctx, &operatorSecret)
				if err != nil {
					return ctrl.Result{}, err
				}

				l.Info("Created secret!")

			} else {
				return ctrl.Result{}, err
			}
		}

		l.Info("Found Secret", "password", string(operatorSecret.Data["password"]))

		// It exists - OK. The configmap for node-red will be updated in a later
		// step, so we do nothing now
		instance.Spec.AdminAuth.Users = append(instance.Spec.AdminAuth.Users, noderednerdendev1.NoderedUser{
			Username:    string(operatorSecret.Data["username"]),
			Password:    string(operatorSecret.Data["password"]),
			Permissions: "*",
		})

	}

	var configmapDirty bool

	// Settings
	{
		configMapName := types.NamespacedName{
			Name:      req.Name + "-nodered" + "-settings",
			Namespace: req.Namespace,
		}

		var current corev1.ConfigMap

		tpl, err := template.New("settings").Parse(settingsTemplate)
		if err != nil {
			return ctrl.Result{}, err
		}

		var b bytes.Buffer

		// Exec template with plain password - it can be used to build a hash which
		// allows comparing the desired and reported state without storing the
		// password.
		if err := tpl.Execute(&b, instance.Spec); err != nil {
			return ctrl.Result{}, err
		}

		// Hash It
		sha := sha256.New()
		sha.Write(b.Bytes())
		hash := sha.Sum(nil)
		b64Hash := base64.URLEncoding.EncodeToString(hash)

		// Now, replace plain passwords with their hashes and Exec tpl again
		for i, _ := range instance.Spec.AdminAuth.Users {
			hashed, err := bcrypt.GenerateFromPassword([]byte(instance.Spec.AdminAuth.Users[i].Password), 8)
			if err != nil {
				panic(err) // TODO
			}
			instance.Spec.AdminAuth.Users[i].Password = string(hashed)
		}

		// Exec tpl again
		b.Reset()
		if err := tpl.Execute(&b, instance.Spec); err != nil {
			return ctrl.Result{}, err
		}

		modified := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name.Name + "-settings",
				Namespace:   req.Namespace,
				Annotations: map[string]string{"nerden.de/data-sha256-sum": b64Hash},
			},
			Data: map[string]string{
				"settings.js": b.String(),
			},
		}

		if err := r.Get(ctx, configMapName, &current); err != nil {
			if errors.IsNotFound(err) {
				err = r.Create(ctx, &modified)
				if err != nil {
					return ctrl.Result{}, err
				}
			} else {
				return ctrl.Result{}, err
			}
		}

		if current.ObjectMeta.Annotations == nil || current.ObjectMeta.Annotations["nerden.de/data-sha256-sum"] != b64Hash {

			l.Info("Need to update cm", "currentSHA256", current.Annotations["nerden.de/data-sha256-sum"], "modifiedSHA256", modified.Annotations["nerden.de/data-sha256-sum"])
			modified.ResourceVersion = current.ResourceVersion
			err := r.Update(ctx, &modified)
			if err != nil {
				return ctrl.Result{}, err
			}

			configmapDirty = true
		}
	}

	// Create sts
	{

		replicas := int32(1)

		image := "nodered/node-red:latest"
		if instance.Spec.Image != "" {
			image = instance.Spec.Image
		}

		modified := appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name.Name,
				Namespace: req.Namespace,
			},
			Spec: appsv1.StatefulSetSpec{
				PodManagementPolicy: appsv1.ParallelPodManagement, // must not be OrderedReady because of https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#forced-rollback
				Replicas:            &replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"nodered.nerden.de/instance": req.Name},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "data",
						},
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("1Gi"),
								},
							},
						},
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"nodered.nerden.de/instance": req.Name},
					},
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{
								Name: "settings",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{Name: name.Name + "-settings"},
									},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:            "nodered",
								Image:           image,
								ImagePullPolicy: corev1.PullIfNotPresent,
								Command: []string{"node-red",
									"-s", "/tmp/settings.js",
									"-u", "/data",
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "data",
										MountPath: "/data",
									},
									{
										Name:      "settings",
										MountPath: "/tmp",
									},
								},
							},
						},
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(instance, &modified, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		var currentSts appsv1.StatefulSet
		if err := r.Get(ctx, name, &currentSts); err != nil {
			if errors.IsNotFound(err) {
				if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(&modified); err != nil {
					return ctrl.Result{}, err
				}

				err = r.Create(ctx, &modified)
				if err != nil {
					return ctrl.Result{}, err
				}
				l.Info("Created StatefulSet", "name", req.Name, "namespace", req.Namespace)

			} else {
				return ctrl.Result{}, err
			}
		}

		if configmapDirty {
			if modified.Spec.Template.ObjectMeta.Annotations == nil {
				modified.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
			}
			modified.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedA"] = time.Now().Format(time.RFC3339)
		}

		patchResult, err := patch.DefaultPatchMaker.Calculate(&currentSts, &modified, compareOpts...)
		if err != nil {
			return ctrl.Result{}, err
		}

		if !patchResult.IsEmpty() {
			l.Info("Update sts", "patch", patchResult.Patch)

			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(&modified); err != nil {
				return ctrl.Result{}, err
			}

			if err := r.Update(ctx, &modified); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Create service
	{
		var current corev1.Service
		modified := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name.Name,
				Namespace: req.Namespace,
			},
			Spec: corev1.ServiceSpec{

				Ports: []corev1.ServicePort{
					{
						Name: "http",
						Port: 1880,
					},
				},
				Selector: map[string]string{"nodered.nerden.de/instance": req.Name},
			},
		}
		if err := controllerutil.SetControllerReference(instance, &modified, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, name, &current); err != nil {
			if errors.IsNotFound(err) {
				if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(&modified); err != nil {
					return ctrl.Result{}, err
				}

				err = r.Create(ctx, &modified)
				if err != nil {
					return ctrl.Result{}, err
				}
			} else {
				return ctrl.Result{}, err
			}
		}

		patchResult, err := patch.DefaultPatchMaker.Calculate(&current, &modified, compareOpts...)
		if err != nil {
			return ctrl.Result{}, err
		}

		if !patchResult.IsEmpty() {
			if err := patch.DefaultAnnotator.SetLastAppliedAnnotation(&modified); err != nil {
				return ctrl.Result{}, err
			}

			modified.ResourceVersion = current.ResourceVersion
			err := r.Update(ctx, &modified)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NoderedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&noderednerdendev1.Nodered{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}
