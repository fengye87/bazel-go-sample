package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	samplev1alpha1 "github.com/fengye87/bazel-go-sample/operator/api/v1alpha1"
)

// GreeterReconciler reconciles a Greeter object
type GreeterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	GreeterServerImage string
	GreeterClientImage string
}

func (r *GreeterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("greeter", req.NamespacedName)

	var greeter samplev1alpha1.Greeter
	if err := r.Client.Get(ctx, req.NamespacedName, &greeter); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get Greeter: %v", err)
	}

	greeterCopy := greeter.DeepCopy()
	if err := r.reconcile(ctx, greeterCopy); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile Greeter: %v", err)
	}

	if !reflect.DeepEqual(greeterCopy.Status, greeter.Status) {
		if err := r.Client.Status().Update(ctx, greeterCopy); err != nil {
			if apierrors.IsConflict(err) {
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, fmt.Errorf("failed to update status of Greeter: %v", err)
		}
		log.Info("updated Greeter status", "status", greeterCopy.Status)
	}

	log.V(1).Info("reconciled Greeter")
	return ctrl.Result{}, nil
}

func (r *GreeterReconciler) reconcile(ctx context.Context, greeter *samplev1alpha1.Greeter) error {
	log := logr.FromContextOrDiscard(ctx)

	if !greeter.DeletionTimestamp.IsZero() {
		log.V(1).Info("Greeter is being deleted")
		return nil
	}

	greeterServerDeploymentName := fmt.Sprintf("%s-greeter-server", greeter.Name)
	greeterServerPodLabels := map[string]string{"app": greeterServerDeploymentName}
	greeterServerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      greeterServerDeploymentName,
			Namespace: greeter.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: greeterServerPodLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: greeterServerPodLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "greeter-server",
						Image: r.GreeterServerImage,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 50051,
						}},
					}},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(greeter, greeterServerDeployment, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller of Deployment %s: %v", greeterServerDeployment.Name, err)
	}
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, greeterServerDeployment, func() error { return nil }); err != nil {
		return fmt.Errorf("failed to create/update Deployment %s: %v", greeterServerDeployment.Name, err)
	}

	greeterServerServiceName := fmt.Sprintf("%s-greeter-server", greeter.Name)
	greeterServerService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      greeterServerServiceName,
			Namespace: greeter.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: greeterServerPodLabels,
			Ports: []corev1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt(50051),
			}},
		},
	}
	if err := controllerutil.SetControllerReference(greeter, greeterServerService, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller of Service %s: %v", greeterServerService.Name, err)
	}
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, greeterServerService, func() error { return nil }); err != nil {
		return fmt.Errorf("failed to create/update Service %s: %v", greeterServerService.Name, err)
	}

	greeterClientDaemonSetName := fmt.Sprintf("%s-greeter-client", greeter.Name)
	greeterClientPodLabels := map[string]string{"app": greeterClientDaemonSetName}
	greeterClientDaemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      greeterClientDaemonSetName,
			Namespace: greeter.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: greeterClientPodLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: greeterClientPodLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:    "greeter-client",
						Image:   r.GreeterClientImage,
						Command: []string{"sleep", "infinity"},
					}},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(greeter, greeterClientDaemonSet, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller of DaemonSet %s: %v", greeterClientDaemonSet.Name, err)
	}
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, greeterClientDaemonSet, func() error { return nil }); err != nil {
		return fmt.Errorf("failed to create/update DaemonSet %s: %v", greeterClientDaemonSet.Name, err)
	}

	return nil
}

func (r *GreeterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplev1alpha1.Greeter{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.DaemonSet{}).
		Complete(r)
}
