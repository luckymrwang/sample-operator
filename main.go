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

package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	inspurincloudv1 "samp/api/v1"
	"samp/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(inspurincloudv1.AddToScheme(scheme))
	utilruntime.Must(inspurincloudv1.AddExtraToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", true,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	//var namespaces []string // List of Namespaces
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		//NewCache:           cache.MultiNamespacedCacheBuilder(namespaces),
		MetricsBindAddress:      "0",
		Port:                    9443,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "c00c35c4.my.domain",
		LeaderElectionNamespace: "kube-system",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.ServiceReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("Service"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("Service"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Service")
		os.Exit(1)
	}
	if err = (&controllers.DeploymentReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("Deployment1"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("Deployment1"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Deployment1")
		os.Exit(1)
	}
	//if err = (&controllers.Deployment2Reconciler{
	//	Client:   mgr.GetClient(),
	//	Log:      ctrl.Log.WithName("controllers").WithName("Deployment2"),
	//	Scheme:   mgr.GetScheme(),
	//	Recorder: mgr.GetEventRecorderFor("Deployment2"),
	//}).SetupWithManager(mgr); err != nil {
	//	setupLog.Error(err, "unable to create controller", "controller", "Deployment2")
	//	os.Exit(1)
	//}

	if err = (&controllers.DestinationRuleReconciler{
		Config:   mgr.GetConfig(),
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("DestinationRule"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("Deployment2"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Deployment2")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
