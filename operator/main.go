package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	samplev1alpha1 "github.com/fengye87/bazel-go-sample/operator/api/v1alpha1"
	"github.com/fengye87/bazel-go-sample/operator/controllers"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = samplev1alpha1.AddToScheme(scheme)
}

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     true,
		LeaderElectionID:   "sample.fengye87.me",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	greeterServerImage := os.Getenv("GREETER_SERVER_IMAGE")
	if greeterServerImage == "" {
		greeterServerImage = "fengye87/greeter_server"
	}
	greeterClientImage := os.Getenv("GREETER_CLIENT_IMAGE")
	if greeterClientImage == "" {
		greeterClientImage = "fengye87/greeter_client"
	}

	if err = (&controllers.GreeterReconciler{
		Client:             mgr.GetClient(),
		Log:                ctrl.Log.WithName("controllers").WithName("Greeter"),
		Scheme:             mgr.GetScheme(),
		GreeterServerImage: greeterServerImage,
		GreeterClientImage: greeterClientImage,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Greeter")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
