package handler

import (
	"reflect"
	"testing"

	opv1 "github.com/kiegroup/kie-cloud-operator/pkg/apis/kiegroup/v1"
	appsv1 "github.com/openshift/api/apps/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/sirupsen/logrus"
)

func TestTrialEnvironmentHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object: &opv1.App{
			Spec: opv1.AppSpec{
				Environment: "trial",
			},
		},
		Deleted: false}
	logrus.Debugf("Testing with environment %v", event.Object.(*opv1.App).Spec.Environment)

	defer func() {
		err := recover().(error)
		logrus.Debugf("Failed with error %v", err)
		assert.Contains(t, err.Error(), "no such host", "Did not get expected no such host error")
	}()

	handler.Handle(nil, event)
}

func TestUnknownEnvironmentHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object: &opv1.App{
			Spec: opv1.AppSpec{
				Environment: "unknown",
			},
		},
		Deleted: false}
	logrus.Debugf("Testing with environment %v", event.Object.(*opv1.App).Spec.Environment)

	defer func() {
		err := recover().(error)
		logrus.Debugf("Failed with error %v", err)
		assert.Contains(t, err.Error(), "invalid memory address or nil pointer dereference", "Did not get expected error")
	}()

	err := handler.Handle(nil, event)
	assert.Nil(t, err, "Unknown environment should result in a no-op and return nil")
}

func TestUnknownResourceTypeHandling(t *testing.T) {
	handler := NewHandler()
	event := sdk.Event{
		Object:  nil,
		Deleted: false}
	logrus.Debugf("Testing with event object %v", reflect.TypeOf(event.Object))

	err := handler.Handle(nil, event)
	assert.Nil(t, err, "Unknown event type should result in a no-op and return nil")
}

func TestUnknownEnvironmentObjects(t *testing.T) {
	cr := &opv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
		},
		Spec: opv1.AppSpec{
			Environment: "unknown",
		},
	}
	logrus.Debugf("Testing with environment %v", cr.Spec.Environment)
	objects, err := NewEnv(cr)
	assert.Equal(t, []runtime.Object{}, objects)
	assert.Equal(t, "envs/unknown.yaml does not exist, environment not deployed", err.Error())
}

func TestTrialEnvironmentObjects(t *testing.T) {
	cr := &opv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
		},
		Spec: opv1.AppSpec{
			Environment: "trial",
		},
	}
	logrus.Debugf("Testing with environment %v", cr.Spec.Environment)
	objects, err := NewEnv(cr)
	assert.Equal(t, "trial-env-rhpamcentr", objects[0].(*appsv1.DeploymentConfig).Name)
	assert.Equal(t, cr.Namespace, objects[0].(*appsv1.DeploymentConfig).ObjectMeta.Namespace)
	assert.Nil(t, err)
}
func TestProdEnvironmentObjects(t *testing.T) {
	cr := &opv1.App{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
		},
		Spec: opv1.AppSpec{
			Environment: "prod",
		},
	}
	logrus.Debugf("Testing with environment %v", cr.Spec.Environment)
	objects, err := NewEnv(cr)
	assert.Equal(t, "prod-rhpamcentr-claim", objects[0].(*corev1.PersistentVolumeClaim).Name)
	assert.Equal(t, cr.Namespace, objects[0].(*corev1.PersistentVolumeClaim).ObjectMeta.Namespace)
	assert.Nil(t, err)
}