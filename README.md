**Warning: This project is in the early stages of prototyping and will be significantly changing over the coming weeks.** 

# Project Goal
To allow EC2 Container Service Task and Service definitions to be used to deploy Pods and Services within Kubernets.

# Prototype Design
The preliminary design uses an aggregated extension API to manage a new ECSDeployment resource type.

The ECSDeployment type encapsulates two sub-resources: ECSTask and ECSService. Both of those sub-resources have the identical structure to their AWS ECS counterparts, and can therefore be copy-pasted into the ECSDeployment definition.

An alternative design would be to reference one or more configmaps from the ECSDeployment and read the Service and Task definitions from those configmaps. Users could then leverage the `kubectl create configmap --from-file` command to avoid having to copy-paste their ECS json/yaml into a ECSDeployment resource. The downside would be a less direct integration with the existing API Machinery. However, these two designs don't have to be mutually exclusive, and the config-map approach may be able to re-use much of the sub-resource approach. Other design proposals are also welcome.

# Installation
The apiserver.yaml is currently the default generated using `make create-agg-config`
```
// Install the API Server and a backing etcd datastore
$ kubectl create -f config/apiserver.yaml
```

# Usage 
These steps can be performed using JSON or YAML formats. The example below uses JSON because it's based on an ECS example:

1. Start by creating the boilerplate for the ECSDeployment resource:
```
{
    "apiVersion": "ecskube.ecskube/v1alpha1",
    "kind": "ECSDeployment",
    "metadata": {
        "labels": {
            "demo": "true"
        },
        "name": "ecs-demo-deploy",
        "namespace": "default"
    },
    "spec": {
        "Task":{}
        "Service:{}
    }
}
```
2. Find the ECS Task definition you're going to deploy. In this case we'll use the wordpress example from https://docs.aws.amazon.com/AmazonECS/latest/developerguide/example_task_definitions.html
```
{
  "containerDefinitions": [
    {
      "name": "wordpress",
      "links": [
        "mysql"
      ],
      "image": "wordpress",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 80,
          "hostPort": 80
        }
      ],
      "memory": 500,
      "cpu": 10
    },
    {
      "environment": [
        {
          "name": "MYSQL_ROOT_PASSWORD",
          "value": "password"
        }
      ],
      "name": "mysql",
      "image": "mysql",
      "cpu": 10,
      "memory": 500,
      "essential": true
    }
  ],
  "family": "hello_world"
}
```
3. Paste the Task into the ECSDeployment's Task property. **Note, while the original Task syntax is valid, the runtime environment is not, so we need to add a couple environment variables to ensure wordpress can find MySQL. This is all done using ECS syntax.**
```
    ...
    "spec": {
        "Task": {
            "Properties": {
                "ContainerDefinitions": [
                    {
                        "Cpu": 100,
                        "Environment": [
                            {
                                "Name": "WORDPRESS_DB_HOST",
                                "Value": "127.0.0.1"
                            },
                            {
                                "Name": "WORDPRESS_DB_PASSWORD",
                                "Value": "password"
                            }
                        ],
                        "Image": "wordpress",
                        "Memory": 500,
                        "Name": "wordpress",
                        "PortMappings": [
                            {
                                "ContainerPort": 80,
                                "HostPort": 80
                            }
                        ]
                    },
                    {
                        "Cpu": 100,
                        "Environment": [
                            {
                                "Name": "MYSQL_ROOT_PASSWORD",
                                "Value": "password"
                            }
                        ],
                        "Essential": true,
                        "Image": "mysql",
                        "Memory": 500,
                        "Name": "mysql"
                    }
                ]
            }
        }
    }
    ...
```

4. Paste an ECS Service Definition as the ECSDeployment's Service property (here we create a minimal one, since the original example doesn't include it).
```
    "spec": {
        "Task": {
           ...
        },
        "Service": {
            "DesiredCount": 1,
            "HealthCheckGracePeriodSeconds": 10,
            "LoadBalancers": [
                {
                    "ContainerName": "wordpress",
                    "ContainerPort": 80
                }
            ],
            "TaskDefinition": "Ignored/implied for now"
        }
        ...
```
See samples/ecsdeployment.json for the full contents.

5. Create the ECS Deployment resource and wait for the service to be exposed
```
$ kubectl create -f samples/ecsdeployment.json
$ kubectl get svc
NAME                 TYPE           CLUSTER-IP       EXTERNAL-IP        PORT(S)        AGE
ecs-demo-deploy      LoadBalancer   100.71.61.71     a8827a69a2e16...   80:30175/TCP   4m

$ kubectl get deployments -l demo=true
NAME              DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
ecs-demo-deploy   1         1         1            1           4m

$ kubectl get pods -l demo=true
NAME                               READY     STATUS    RESTARTS   AGE
ecs-demo-deploy-568d58fb5c-lq4fs   2/2       Running   0          5m

```

Once DNS for the new service has propagated, browsing to port 80 should show the wordpress setup wizard.

6. Delete the ECS Deployment and watch the other resources terminate
```
$ kubectl delete ecsdeployment ecs-demo-deploy
ecsdeployment "ecs-demo-deploy" deleted
$ kubectl get deployment,services,pods -l demo=true
No resources found.
```


Note: update is not supported yet
# Development
Note: this project needs to be deployed into k8s 1.8 or later.

## Project structure
We use the extremely helpful [apiserver-builder tool](https://github.com/kubernetes-incubator/apiserver-builder) to generate the initial project structure and manage the API machinery.

* **cmd/**: contains the apiserver and controller-manager entry points. Both of these were generated by apiserver-builder.

* **config/**: contains the yaml for deploying the ecs-kube-apisever and backing api store, as well as required certificates. This was generated by apiserver-builder using `make build-agg-config`

* **pkg/**: contains a mix of generated and custom api/controller code - in essentially the standard layout you'll find in all K8s api code. The generated code can be regenerated using `make build-generated`

## Crash Course on apiserver-builder:
TLDR; To add fields or resources, edit pkg/apis/ecskube/v1alpha1/ecskube_types.go then regenerate the API machinery. To add business logic, edit pkg/controller/ecskube/controller.go.

This is only meant to provide minimal context for understanding this codebase. You'll  want to read the apiserver-builder docs before making any significant changes.

The two most important functions of this code base are: 
1. The definition of custom resources
1. The logic for handling those custom resources

### Custom Resource Definitions
These live in pkg/apis/ecskube/v1alpha1/ecskube_types.go. The resources are a bunch of structs that correspond to a particular version (in this case v1alpha1) of the API. Any changes to these structs need to be reflected in the generated API machinery. Running `apiserver-boot build generated` reads the resources definitions, rebuilds the API machinery, and creates an unversioned (ie internal and canonical) representation of each resource.

### Custom Logic
Business logic is handled in pkg/controller/ecskube/controller.go and the controller methods are responsible for watching and reacting to changes in our custom resources. The two important methods are:
* **Init**: Sets up watches on the K8s resources we care about and stores references to the corresponding listers in the controller
* **Reconcile**: Gets called when custom resources are created are updated. Handles all the heavy lifting of converting the custom resources into Services, Deployments, etc and then updating those representations via the main API. In other words, this is where we create/update things

#### Garbage Collection
You'll notice there's no "delete" method mentioned above.  That's because we rely on K8s Garbage Collection to delete Services, Deployment, etc when the corresponding ECSDeployment is deleted. All that's needed is to ensure we set the OwnerReferences to point to the ECSDeployment on any resources we create. 


### Commands Used to Create this Repo
For posterity...
```
$ apiserver-boot init repo --domain ecskube
$ apiserver-boot create group version resource --group ecskube --version v1alpha1 --kind ECSDeployment
```

# Contributing
Contributions are always welcome. Especially if you can make sense of any of the above at this stage ;)