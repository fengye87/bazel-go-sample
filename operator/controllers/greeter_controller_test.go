package controllers

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	samplev1alpha1 "github.com/fengye87/bazel-go-sample/operator/api/v1alpha1"
)

var _ = Describe("GreeterController", func() {
	Context("for a Greeter object", func() {
		var greeterKey types.NamespacedName

		BeforeEach(func() {
			greeter := &samplev1alpha1.Greeter{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "t-",
					Namespace:    "default",
				},
				Spec: samplev1alpha1.GreeterSpec{},
			}
			err := k8sClient.Create(context.Background(), greeter)
			Expect(err).ToNot(HaveOccurred())
			greeterKey = types.NamespacedName{
				Name:      greeter.Name,
				Namespace: greeter.Namespace,
			}
		})

		It("should create a greeter-server Deployment", func() {
			Eventually(func() *appsv1.Deployment {
				var greeterServerDeployment appsv1.Deployment
				greeterServerDeploymentKey := types.NamespacedName{
					Name:      fmt.Sprintf("%s-greeter-server", greeterKey.Name),
					Namespace: greeterKey.Namespace,
				}
				err := k8sClient.Get(context.Background(), greeterServerDeploymentKey, &greeterServerDeployment)
				if apierrors.IsNotFound(err) {
					return nil
				}
				Expect(err).ToNot(HaveOccurred())
				return &greeterServerDeployment
			}).ShouldNot(BeNil())
		})

		It("should create a greeter-server Service", func() {
			Eventually(func() *corev1.Service {
				var greeterServerService corev1.Service
				greeterServerServiceKey := types.NamespacedName{
					Name:      fmt.Sprintf("%s-greeter-server", greeterKey.Name),
					Namespace: greeterKey.Namespace,
				}
				err := k8sClient.Get(context.Background(), greeterServerServiceKey, &greeterServerService)
				if apierrors.IsNotFound(err) {
					return nil
				}
				Expect(err).ToNot(HaveOccurred())
				return &greeterServerService
			}).ShouldNot(BeNil())
		})

		It("should create a greeter-client DaemonSet", func() {
			Eventually(func() *appsv1.DaemonSet {
				var greeterClientDaemonSet appsv1.DaemonSet
				greeterClientDaemonSetKey := types.NamespacedName{
					Name:      fmt.Sprintf("%s-greeter-client", greeterKey.Name),
					Namespace: greeterKey.Namespace,
				}
				err := k8sClient.Get(context.Background(), greeterClientDaemonSetKey, &greeterClientDaemonSet)
				if apierrors.IsNotFound(err) {
					return nil
				}
				Expect(err).ToNot(HaveOccurred())
				return &greeterClientDaemonSet
			}).ShouldNot(BeNil())
		})
	})
})
