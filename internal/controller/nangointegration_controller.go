/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	nangov1alpha1 "github.com/rossmcewan/nango-integration-operator/api/v1alpha1"
	"github.com/rossmcewan/nango-integration-operator/internal/nango"
)

// NangoIntegrationReconciler reconciles a NangoIntegration object
type NangoIntegrationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nango.nango.dev,resources=nangointegrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nango.nango.dev,resources=nangointegrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nango.nango.dev,resources=nangointegrations/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NangoIntegration object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *NangoIntegrationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Starting reconciliation", "namespace", req.Namespace, "name", req.Name)

	// Fetch the NangoIntegration instance
	nangoIntegration := &nangov1alpha1.NangoIntegration{}
	err := r.Get(ctx, req.NamespacedName, nangoIntegration)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Info("NangoIntegration resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get NangoIntegration")
		return ctrl.Result{}, err
	}

	// Resolve secret references
	clientID, err := r.resolveSecretOrStringValue(ctx, nangoIntegration.Namespace, nangoIntegration.Spec.Credentials.ClientID)
	if err != nil {
		log.Error(err, "Failed to resolve client ID")
		return r.updateStatus(ctx, nangoIntegration, "Failed", fmt.Sprintf("Failed to resolve client ID: %v", err), err)
	}

	clientSecret, err := r.resolveSecretOrStringValue(ctx, nangoIntegration.Namespace, nangoIntegration.Spec.Credentials.ClientSecret)
	if err != nil {
		log.Error(err, "Failed to resolve client secret")
		return r.updateStatus(ctx, nangoIntegration, "Failed", fmt.Sprintf("Failed to resolve client secret: %v", err), err)
	}

	nangoToken, err := r.resolveSecretOrStringValue(ctx, nangoIntegration.Namespace, nangoIntegration.Spec.NangoToken)
	if err != nil {
		log.Error(err, "Failed to resolve Nango token")
		return r.updateStatus(ctx, nangoIntegration, "Failed", fmt.Sprintf("Failed to resolve Nango token: %v", err), err)
	}

	// Check if the integration already exists in Nango
	nangoClient := nango.NewClient(nangoIntegration.Spec.NangoBaseURL, nangoToken)

	// Try to get the integration from Nango
	_, err = nangoClient.GetIntegration(nangoIntegration.Spec.UniqueKey)
	if err == nil {
		// Integration exists, update status
		log.Info("Integration already exists in Nango", "uniqueKey", nangoIntegration.Spec.UniqueKey)
		return r.updateStatus(ctx, nangoIntegration, "Created", "", nil)
	}

	// Integration doesn't exist, create it
	log.Info("Creating integration in Nango", "uniqueKey", nangoIntegration.Spec.UniqueKey)

	// Prepare the request
	createReq := nango.CreateIntegrationRequest{
		UniqueKey:   nangoIntegration.Spec.UniqueKey,
		Provider:    nangoIntegration.Spec.Provider,
		DisplayName: nangoIntegration.Spec.DisplayName,
		Credentials: nango.NangoCredentials{
			Type:         nangoIntegration.Spec.Credentials.Type,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       nangoIntegration.Spec.Credentials.Scopes,
		},
	}

	// Create the integration
	response, err := nangoClient.CreateIntegration(createReq)
	if err != nil {
		log.Error(err, "Failed to create integration in Nango")
		return r.updateStatus(ctx, nangoIntegration, "Failed", err.Error(), err)
	}

	log.Info("Successfully created integration in Nango",
		"uniqueKey", response.Data.UniqueKey,
		"provider", response.Data.Provider)

	return r.updateStatus(ctx, nangoIntegration, "Created", "", nil)
}

// resolveSecretOrStringValue resolves a value from either a direct string or a secret reference
func (r *NangoIntegrationReconciler) resolveSecretOrStringValue(ctx context.Context, namespace string, secretOrString nangov1alpha1.SecretOrStringValue) (string, error) {
	// If direct value is provided, use it
	if secretOrString.Value != "" {
		return secretOrString.Value, nil
	}

	// If secret reference is provided, resolve it
	if secretOrString.SecretKeyRef != nil {
		secret := &corev1.Secret{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      secretOrString.SecretKeyRef.Name,
		}, secret)
		if err != nil {
			return "", fmt.Errorf("failed to get secret %s: %w", secretOrString.SecretKeyRef.Name, err)
		}

		value, exists := secret.Data[secretOrString.SecretKeyRef.Key]
		if !exists {
			return "", fmt.Errorf("key %s not found in secret %s", secretOrString.SecretKeyRef.Key, secretOrString.SecretKeyRef.Name)
		}

		return string(value), nil
	}

	return "", fmt.Errorf("neither value nor secretKeyRef provided")
}

// updateStatus updates the status of the NangoIntegration resource
func (r *NangoIntegrationReconciler) updateStatus(ctx context.Context, nangoIntegration *nangov1alpha1.NangoIntegration, status, errorMessage string, err error) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Update the status
	nangoIntegration.Status.Status = status
	nangoIntegration.Status.ErrorMessage = errorMessage
	nangoIntegration.Status.LastUpdated = &metav1.Time{Time: time.Now()}

	// Add condition
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "IntegrationCreated",
		Message:            "Integration successfully created in Nango",
		LastTransitionTime: metav1.Time{Time: time.Now()},
	}

	if status == "Failed" {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "IntegrationCreationFailed"
		condition.Message = fmt.Sprintf("Failed to create integration: %s", errorMessage)
	}

	// Update or add the condition
	conditionUpdated := false
	for i, existingCondition := range nangoIntegration.Status.Conditions {
		if existingCondition.Type == condition.Type {
			nangoIntegration.Status.Conditions[i] = condition
			conditionUpdated = true
			break
		}
	}
	if !conditionUpdated {
		nangoIntegration.Status.Conditions = append(nangoIntegration.Status.Conditions, condition)
	}

	// Update the status in Kubernetes
	if err := r.Status().Update(ctx, nangoIntegration); err != nil {
		log.Error(err, "Failed to update NangoIntegration status")
		return ctrl.Result{}, err
	}

	if status == "Failed" {
		// If creation failed, requeue after a delay
		return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NangoIntegrationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nangov1alpha1.NangoIntegration{}).
		Named("nangointegration").
		Complete(r)
}
