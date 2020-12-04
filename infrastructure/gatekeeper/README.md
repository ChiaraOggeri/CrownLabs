

GATEKEEPER
=====================================
## GOALS 
Kubernetes allows you to write policies on your cluster objects by means of admission controller webhooks which are executed everytime a cluster component is created or modified.
OPA (Open Policy Agent) is an open source, general-purpose policy engine to write unifed policies across different applications. It exploits REGO language to write policies.
Gatekeeper allows you to integrate easily OPA on kubernetes by adding:
 - native kubernetes CRDs for instantiating the policy library (Constraints)
 - native kubernetes CRDs for extending the policy library (ConstraintTemplates)


## HOW TO INSTALL
### Prerequisites 
- minimun kubernetes version: 1.14
- make sure you have cluster admin permissions
#### Deploying a Release using Prebuilt Image
Run the following command:

kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/release-3.1/deploy/gatekeeper.y

## HOW TO USE GATEKEEPER
### ConstraintTemplates
In the constraintemplates you should define your own constraint CRD and associate it to your rego policies, for example :
   **constrainttemplate.yaml**

     apiVersion: templates.gatekeeper.sh/v1beta1
     kind: ConstraintTemplate
     metadata:
       name: k8srequiredlabels
     spec:
       crd: #definition of your Constraint CRD
         spec:
           names:
             kind: K8sRequiredLabels
           validation:
              # Schema for the `parameters` field
            openAPIV3Schema:
              properties: #if you need some properties define them here
                labels:
                  type: array
                  items: string
       targets:
          - target: admission.k8s.gatekeeper.sh
            rego: |   #define your rego policies here
            # if you need to take some parameters from your constraint specify the package name here
             package k8srequiredlabels 
             #violation is the function that allows you define your policy.
             #It returns two values: msg and details
             violation[{"msg": msg, "details": {"missing_labels": missing}}] {
             #input.review is the object under controll
            provided := {label |input.review.object.metadata.labels[label]}
             #input.parameters are your parameters defined above
            required := {label | label := input.parameters.labels[_]}
            missing := required - provided
            #if the following statement is true msg is evaluated and therefore there is a violation, otherwise your changes on cluster are accepted
            count(missing) > 0
            msg := sprintf("you must provide labels: %v", [missing])
     }
  
 ### Constraint 
In constraint you should declare on which resources your constraintemplate must be enforced and define your parameters values.
**constraint.yaml**

       apiVersion: constraints.gatekeeper.sh/v1beta1
      kind: K8sRequiredLabels
      metadata:
        name: ns-must-have-gk
      spec:
        match:
          kinds:
          # declare your resources to be put under control here
          - apiGroups: [""]
            kinds: ["Namespace"]
          # declare your parameters values here
          parameters:
            labels: ["gatekeeper"]

### Config-Sync
If you need to get resources from the cluster to write your policy yuo need  first to cache those and then you will be able to retrieve them in the ConstraintTemplate rules.
For example:

    apiVersion: config.gatekeeper.sh/v1alpha1
    kind: Config
    metadata:
       name: config
       namespace: "gatekeeper-system"
    spec:
      sync:
        syncOnly:
        # define the objects that you need to cache here
         - group: ""
           version: "v1"
           kind: "Namespace"
         - group: ""
           version: "v1"
           kind: "Pod"
To retrieve the cached resources in the rego rules you can do as follow:

-   For cluster-scoped objects:  `data.inventory.cluster[<groupVersion>][<kind>][<name>]`
    -   Example referencing the Gatekeeper namespace:  `data.inventory.cluster["v1"].Namespace["gatekeeper"]`
-   For namespace-scoped objects:  `data.inventory.namespace[<namespace>][groupVersion][<kind>][<name>]`
    -   Example referencing the Gatekeeper pod:  `data.inventory.namespace["gatekeeper"]["v1"]["Pod"]["gatekeeper-controller-manager-d4c98b788-j7d92"]`
