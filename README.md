# Energy-aware K8S Function Scheduler on Serverless

## Introduction
An energy-aware custom <a href="https://github.com/kubernetes/kubernetes">Kubernetes</a> scheduler project for Serverless environment with purpose of increasing the availability of the services. The project is implemented on <a href="https://github.com/kubernetes-sigs/scheduler-plugins">The Scheduling Framework</a> built by the K8S team and the functions are deployed on <a href="https://github.com/openfaas">OpenFaas</a> which is a serverless deployment tool.     
Although the title highlights the scheduling process but there are other components along it. Each component is documented with details.


Install minkube, kuberntees, openfaas

## System Architecture
<p align="center">
  <img style="display:block;margin-left: auto; margin-right: auto; width:50 %" width="500px" src="system_architecture.png">
</p>

There are three main components on the master node each with a specific responsibility. Components are as follow:

### Energy Adapter    

The Energy Adapter is responsible for fetching the energy level of each node the physical devices, convert the values into discrete values and assign those values to each corresponding node. the energy adapter is being executed periodically, for example every one minute to have the system updates.   
Although this is the main idea behind the base paper but this project is implemented with more abstractions so there's no physical energy meter. The discrete values are generated randomly for each node. The energy levels are as follow:
- well powered: 75 to 100 percent energy
- vulnerable:   40 to 75 percent energy
- low powered:  15 to 40 percent energy
- powerless:    0 to 15 percent energy

These energy levels are being bound to every node using the pods labels. The key for the label is `energy-status`.
This component is placed in the <a href="https://github.com/faramarzQ/energy-aware-k8s-serverless-scheduler/tree/main/energy_adapter">energy_adapter</a> directory. On the first run of the plugin, it assigns new values to each node. on the later runs it switches the energy levels; `well-powered` and `low-powered` labels are switched and `vulnerable` and `powerles` labels too.

### Function Offloader

The Function Offloader component which is also being executed periodically is responsible to keeping the pods on nodes with higher energy levels. For each pod it checks if there is better nodes out there with higher energy levels and empty capacity to host the pod. If a better node found, the pod is being terminated so that the Deployment resource recreates it. This component is placed in the <a href="https://github.com/faramarzQ/energy-aware-k8s-serverless-scheduler/tree/main/function_offloader">function_offloader</a>.

### Function Scheduler

The main component of the system is The Scheduler Component with is called on the creation of each pod and finds the best node for hosting it. This component uses the scheduling framework so in order to understand it you need to take a look at <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/">this documentation</a>. This scheduler only implements the Filter Plugin to filter the nodes to find the best node. This component is placed in the <a href="https://github.com/faramarzQ/energy-aware-k8s-serverless-scheduler/tree/main/function_scheduler">function_scheduler</a>.

## Code Components

There are two components used by the project, the cluster and the clientset.
The clientset creates an instance of the Kubernetes' Clientset module to communicate with the cluster. The cluster component has the overridden models of the Kubernetes Resources including Pod and Node. Each model is responsible for the logic of the corresponding domain like deleting pods of getting nodes' energy status. Packages are placed in the <a href="https://github.com/faramarzQ/energy-aware-k8s-serverless-scheduler/tree/main/pkg">pkg</a>.

## The Function
I've used a simple fibonacci function as the service provided to users which returns the fibonacci value of a given number. this function is highly cpu intensive for large numbers