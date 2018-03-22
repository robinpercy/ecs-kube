/*
Copyright 2017 The Kubernetes Authors.

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

package ecsdeployment

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/robinpercy/ecs-kube/pkg/apis/ecskube/v1alpha1"
	listers "github.com/robinpercy/ecs-kube/pkg/client/listers_generated/ecskube/v1alpha1"
	appsv1beta2listers "k8s.io/client-go/listers/apps/v1beta2"
	corev1listers "k8s.io/client-go/listers/core/v1"

	"github.com/robinpercy/ecs-kube/pkg/controller/sharedinformers"
)

const (
	DEPLOYMENT_LABEL = "ecsDeployment"
	OWNERREF_KIND    = "ECSDeployment"
	OWNERREF_VERSION = "ecskube.robinpercy.com/v1alpha1"
)

// +controller:group=ecskube,version=v1alpha1,kind=ECSDeployment,resource=ecsdeployments
type ECSDeploymentControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about ECSDeployment
	lister           listers.ECSDeploymentLister
	deploymentLister appsv1beta2listers.DeploymentLister
	serviceLister    corev1listers.ServiceLister
	si               *sharedinformers.SharedInformers
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *ECSDeploymentControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	c.si = arguments.GetSharedInformers()
	c.lister = c.si.Factory.Ecskube().V1alpha1().ECSDeployments().Lister()
	deployments := c.si.KubernetesFactory.Apps().V1beta2().Deployments()
	c.deploymentLister = deployments.Lister()
	services := c.si.KubernetesFactory.Core().V1().Services()
	c.serviceLister = services.Lister()
	arguments.Watch("Deployment", deployments.Informer(), c.DeploymentToECSDeploymentId)
	arguments.Watch("Service", services.Informer(), c.ServiceToECSDeploymentId)

}

// Reconcile handles enqueued messages
func (c *ECSDeploymentControllerImpl) Reconcile(u *v1alpha1.ECSDeployment) error {
	// Implement controller logic here
	// TODO: create services too
	// TODO: ensure we don't have a race condition if we deploy multiple controllers
	log.Printf("Running reconcile ECSDeployment for: %s\n", u.Name)
	d := c.si.KubernetesClientSet.AppsV1beta2().Deployments(u.Namespace)

	dList, err := d.List(metav1.ListOptions{
		LabelSelector: DEPLOYMENT_LABEL + "=" + u.GetName(),
	})

	if err != nil {
		err = fmt.Errorf("Failed to lookup deployments for %s: %v", u.Name, err)
		log.Printf("%v", err)
		return err
	}

	if len(dList.Items) == 0 {
		log.Printf("Creating deployment for ECSDeployment  %v", u)
		newDep, err := ecsDeploymentToDeployment(u)
		if err != nil {
			log.Printf("warning: deployment was not created:%v", err)
		}

		log.Printf("Creating deployment  %v", newDep)
		_, err = d.Create(newDep)
		if err != nil {
			log.Printf("%v", err)
			return fmt.Errorf("Failed to create deployment for ECSDeployment '%s': %v", u.GetName(), err)
		}
	}

	s := c.si.KubernetesClientSet.CoreV1().Services(u.Namespace)

	sList, err := s.List(metav1.ListOptions{
		LabelSelector: DEPLOYMENT_LABEL + "=" + u.GetName(),
	})

	if err != nil {
		err = fmt.Errorf("Failed to lookup services for %s: %v", u.Name, err)
		log.Printf("%v", err)
		return err
	}

	if len(sList.Items) == 0 {
		log.Printf("Creating service for ECSDeployment  %v", u)
		newSvc, err := ecsDeploymentToService(u)
		if err != nil {
			log.Printf("warning: service was not created:%v", err)
		}

		log.Printf("Creating service  %v", newSvc)
		_, err = s.Create(newSvc)
		if err != nil {
			log.Printf("%v", err)
			return fmt.Errorf("Failed to create service for ECSDeployment '%s': %v", u.GetName(), err)
		}
	}
	return nil
}

// DeploymentToECSDeploymentId returns a ECSDeployment ID based on a deployment
func (c *ECSDeploymentControllerImpl) DeploymentToECSDeploymentId(i interface{}) (string, error) {
	d, _ := i.(*appsv1beta2.Deployment)
	log.Printf("Deployment update: %v", d.Name)
	if len(d.OwnerReferences) == 1 && d.OwnerReferences[0].Kind == OWNERREF_KIND {
		return d.Namespace + "/" + d.OwnerReferences[0].Name, nil
	} else {
		return "", nil
	}
}

// ServiceToECSDeploymentId returns a ECSDeployment ID based on a service
func (c *ECSDeploymentControllerImpl) ServiceToECSDeploymentId(i interface{}) (string, error) {
	s, _ := i.(*v1.Service)
	log.Printf("Service update: %v", s.Name)
	if len(s.OwnerReferences) == 1 && s.OwnerReferences[0].Kind == OWNERREF_KIND {
		return s.Namespace + "/" + s.OwnerReferences[0].Name, nil
	} else {
		return "", nil
	}
}

func (c *ECSDeploymentControllerImpl) Get(namespace, name string) (*v1alpha1.ECSDeployment, error) {
	return c.lister.ECSDeployments(namespace).Get(name)
}

func ecsDeploymentToService(u *v1alpha1.ECSDeployment) (*v1.Service, error) {

	// TODO: handle multiple (and maybe 0?) LBs
	if len(u.Spec.Service.LoadBalancers) != 1 {
		return nil, fmt.Errorf("Service was not created for ECSDeployment: '%s' since no LoadBalancer config was specified.", u.Name)
	}

	labels := labelsForECSDeployment(u)

	// TODO: This is a bare-minimum implementation.
	// TODO: respect ECS port semantics.
	ports := make([]v1.ServicePort, 0, 2)
	for _, l := range u.Spec.Service.LoadBalancers {
		ports = append(ports, v1.ServicePort{
			Port:       l.ContainerPort, // TODO: figure out if ECS lets you control this via svc definition
			TargetPort: intstr.FromInt(int(l.ContainerPort)),
		})
	}

	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "core/v1",
		},
		ObjectMeta: *objectMetaForECSDeployment(u, labels),
		Spec: v1.ServiceSpec{
			Selector: labels,
			Type:     v1.ServiceTypeLoadBalancer, // TODO: hardcoded for prototype
			Ports:    ports,
		},
	}, nil

}

func ecsDeploymentToDeployment(u *v1alpha1.ECSDeployment) (*appsv1beta2.Deployment, error) {

	containers, err := convertContainers(u.Spec.Task.Properties.ContainerDefinitions)
	if err != nil {
		return nil, fmt.Errorf("Failed to build deployment from ECSDeployment '%s':%v", u.GetName(), err)
	}

	labels := labelsForECSDeployment(u)

	d := &appsv1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: *objectMetaForECSDeployment(u, labels),
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &u.Spec.Service.DesiredCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: containers,
				},
			},
		},
	}

	return d, nil
}

func objectMetaForECSDeployment(u *v1alpha1.ECSDeployment, labels map[string]string) *metav1.ObjectMeta {
	t := true
	return &metav1.ObjectMeta{
		Namespace: u.GetNamespace(),
		Name:      u.GetName(),
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			metav1.OwnerReference{
				UID:                u.GetUID(),
				Name:               u.GetName(),
				Controller:         &t,
				BlockOwnerDeletion: &t,
				Kind:               OWNERREF_KIND,
				APIVersion:         OWNERREF_VERSION,
			},
		},
	}
}

func labelsForECSDeployment(u *v1alpha1.ECSDeployment) map[string]string {
	labels := make(map[string]string)
	for k, v := range u.Labels {
		labels[k] = v
	}
	labels[DEPLOYMENT_LABEL] = u.GetName()
	return labels
}

func convertContainers(cList []v1alpha1.ECSContainer) ([]v1.Container, error) {
	cDefs := make([]v1.Container, len(cList))
	for i, c := range cList {
		rr, err := buildResourceRequirements(c.CPU, c.Memory)
		if err != nil {
			return cDefs, fmt.Errorf("failed to build resource requirements for container %s: %v", c.Name, err)
		}
		cDefs[i] = v1.Container{
			Name:      c.Name,
			Image:     c.Image,
			Ports:     convertPorts(c.PortMappings),
			Env:       convertEnv(c.Environment),
			Resources: rr,
		}
	}

	return cDefs, nil
}

func convertPorts(pList []v1alpha1.ECSPortMapping) []v1.ContainerPort {
	// TODO: figure out how to reconcile the differences in host port usage between ECS and K8s
	// TODO: is it safe to pass the porotocol straight through?
	kports := make([]v1.ContainerPort, len(pList))
	for i, p := range pList {
		proto := v1.ProtocolTCP

		switch strings.ToLower(p.Protocol) {
		case "tcp":
		case "udp":
			proto = v1.ProtocolUDP
		default:
			log.Printf("Invalid protocol '%s' found, defaulting to tcp")
		}

		kports[i] = v1.ContainerPort{
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      proto,
		}
	}
	return kports
}

func convertEnv(kvList []v1alpha1.ECSKeyValuePair) []v1.EnvVar {
	kvars := make([]v1.EnvVar, len(kvList))
	for i, kv := range kvList {
		kvars[i] = v1.EnvVar{
			Name:  kv.Name,
			Value: kv.Value,
		}
	}

	return kvars
}

func buildResourceRequirements(cpu, memory int) (v1.ResourceRequirements, error) {
	// TODO: provide some global config to determine QoS (request vs limit)?
	// TODO: reconcile ECS (relative) cpu semantics with K8s (absolute) semantics, for now, assuming absolute
	// TODO: this needs to be made much more robust
	rr := v1.ResourceRequirements{}
	cpuStr := strconv.Itoa(cpu) + "m"
	memStr := strconv.Itoa(memory) + "M"
	parsedCPU, err := resource.ParseQuantity(cpuStr)
	if err != nil {
		return rr, fmt.Errorf("error parsing value for CPU, '%s':%v", cpuStr, err)
	}
	parsedMem, err := resource.ParseQuantity(memStr)
	if err != nil {
		return rr, fmt.Errorf("error parsing value for Memory, '%s':%v'", memStr, err)
	}
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    parsedCPU,
			v1.ResourceMemory: parsedMem,
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    parsedCPU,
			v1.ResourceMemory: parsedMem,
		},
	}, nil
}
