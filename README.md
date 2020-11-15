# Overview

kubernetes 允许用户自定义自己的资源对象， 自定义资源 `CRD（Custom Resource Definition）` 可以扩展 Kubernetes API， 然后为自定义资源写一个对应的控制器，推出自己的声明式 API。

## CRD

### 什么是 CRD？

CRD 本身是一种 Kubernetes 内置的资源类型，是 CustomResourceDefinition 的缩写，可以通过 `kubectl get crd` 命令查看集群内定义的 CRD 资源。

- 在 Kubernetes，所有的东西都叫做资源（Resource），就是 Yaml 里 Kind 那项所描述的

- 但除了常见的 Deployment 之类的内置资源之外，Kube 允许用户自定义资源（Custom Resource），也就是 CR

- CRD 其实并不是自定义资源，而是我们自定义资源的定义（来描述我们定义的资源是什么样子）


### CRD 能做什么？

一般情况下，我们利用 CRD 所定义的 CR 就是一个新的控制器，我们可以自定义控制器的逻辑，来做一些 Kubernetes 集群原生不支持的功能。

- 利用 CRD 所定义的 CR 就是一个新的控制器，我们可以自定义控制器的逻辑，来做一些 Kubernetes 集群原生不支持的功能。
- CRD 使得 Kubernetes 已有的资源和能力变成了乐高积木，我们很轻松就可以利用这些积木拓展 Kubernetes 原生不具备的能力。
- 其次是产品上，基于 Kubernetes 做的产品无法避免的需要让我们将产品术语向 Kube 术语靠拢，比如一个服务就是一个 Deployment，一个实例就是一个 Pod 之类。但是 CRD 允许我们自己基于产品创建概念（或者说资源），让 Kube 已有的资源为我们的概念服务，这可以使产品更专注与解决的场景，而不是如何思考如何将场景应用到 Kubernetes。
- CRD 允许我们基于已有的 Kubernetes 资源，拓展集群能力
- CRD 可以使我们自己定义一套成体系的规范，自造概念

### 怎么实现 CRD 扩展？

- 编写 CRD 并将其部署到 Kubernetes 集群里；

   这一步的作用就是让 Kubernetes 知道有这个资源及其结构属性，在用户提交该自定义资源的定义时（通常是 YAML 文件定义），Kubernetes 能够成功校验该资源并创建出对应的 Go struct 进行持久化，同时触发控制器的调谐逻辑。

- 编写 Controller 并将其部署到 Kubernetes 集群里。


## Kubebuilder

Kubebuilder 节省大量工作，方便用户从零开始开发 CRDs，Controllers 和 Admission Webhooks，让扩展 Kubernetes 变得更简单

### Installation

```bash
cat <<EOF | tee ./script/install_kubebuilder.sh
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
EOF
```

```bash
chmod +x ./script/install_kubebuilder.sh
./script/install_kubebuilder.sh
```

## Step-by-step write CR

1. 初始化项目
    
    ```bash
    # 创建了一个 Go module 工程，同时创建了一些模板文件。
    kubebuilder init --domain basebit.me --repo github.com/lqshow/kubernetes-crd --owner "LQ"
    ```
	
	**Output**
	```bash
	Writing scaffold for you to edit...
	Get controller runtime:
	$ go get sigs.k8s.io/controller-runtime@v0.5.0
	Update go.mod:
	$ go mod tidy
	Running make:
	$ make
	/Users/linqiong/workspace/app/golang/lib/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
	go fmt ./...
	go vet ./...
	go build -o bin/manager main.go
	Next: define a resource with:
	$ kubebuilder create api
	```
	
	**脚手架工程结构**
	```bash
	➜ tree
	.
	├── Dockerfile # Build Controller 镜像
	├── Makefile # 用于构建和部署 controller
	├── PROJECT # 用于创建新组件的 Kubebuilder 元数据
	├── bin
	│   └── manager # 编译后的可执行二进制文件
	├── config
	│   ├── certmanager
	│   │   ├── certificate.yaml
	│   │   ├── kustomization.yaml
	│   │   └── kustomizeconfig.yaml
	│   ├── default
	│   │   ├── kustomization.yaml
	│   │   ├── manager_auth_proxy_patch.yaml
	│   │   ├── manager_webhook_patch.yaml
	│   │   └── webhookcainjection_patch.yaml
	│   ├── manager # 部署 Controller 的 manifest 文件
	│   │   ├── kustomization.yaml
	│   │   └── manager.yaml
	│   ├── prometheus
	│   │   ├── kustomization.yaml
	│   │   └── monitor.yaml
	│   ├── rbac # 运行 Controller 需要的 RBAC 权限
	│   │   ├── auth_proxy_client_clusterrole.yaml
	│   │   ├── auth_proxy_role.yaml
	│   │   ├── auth_proxy_role_binding.yaml
	│   │   ├── auth_proxy_service.yaml
	│   │   ├── kustomization.yaml
	│   │   ├── leader_election_role.yaml
	│   │   ├── leader_election_role_binding.yaml
	│   │   └── role_binding.yaml
	│   └── webhook
	│       ├── kustomization.yaml
	│       ├── kustomizeconfig.yaml
	│       └── service.yaml
	├── go.mod
	├── go.sum
	├── hack
	│   └── boilerplate.go.txt
	└── main.go # Controller Entrypoint
	```
    
2. 创建新 API
    
    ```bash
    # 创建 API 后，kubebuilder 会创建 crd 数据定义文件以及对应的 controller 文件
    kubebuilder create api --group runner --version v1alpha1 --kind App
    kubebuilder create api --group runner --version v1alpha1 --kind Fuwu
    ```
	
	**Output**
	```bash
	Create Resource [y/n]
	y
	Create Controller [y/n]
	y
	Writing scaffold for you to edit...
	api/v1alpha1/app_types.go
	controllers/app_controller.go
	Running make:
	$ make
	/Users/linqiong/workspace/app/golang/lib/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
	go fmt ./...
	go vet ./...
	go build -o bin/manager main.go
	```
	
	**需要开发者关注的增量目录**
	```bash
	➜ tree  -I "hack|bin|webhook|rbac|prometheus|manager|default|certmanager|main.go|go.mod|go.sum|*file|PROJECT"
	.
	├── api
	│   └── v1alpha1 # 自定义 CRD 目录（xx_types.go）
	│       ├── app_types.go
	│       ├── fuwu_types.go
	│       ├── groupversion_info.go
	│       └── zz_generated.deepcopy.go
	├── config
	│   ├── crd # 部署 CRD 的 manifest 文件
	│   │   ├── kustomization.yaml
	│   │   ├── kustomizeconfig.yaml
	│   │   └── patches
	│   │       ├── cainjection_in_apps.yaml
	│   │       ├── cainjection_in_fuwus.yaml
	│   │       ├── webhook_in_apps.yaml
	│   │       └── webhook_in_fuwus.yaml
	│   └── samples # CRD 样例 manifest 文件
	│       ├── runner_v1alpha1_app.yaml
	│       └── runner_v1alpha1_fuwu.yaml
	└── controllers # 自定义 Controller 逻辑
		├── app_controller.go
		├── fuwu_controller.go
		└── suite_test.go
	```
    
    **查看生成的 sample yaml 文件**
    ```yaml
    # cat config/samples/runner_v1alpha1_app.yaml
    apiVersion: runner.basebit.me/v1alpha1
    kind: App
    metadata:
      name: app-sample
    spec:
      # Add fields here
      foo: bar
    ```
    
3. 定义 CRD
    
    **主要关注 2 个文件**
    -  CRD 的定义文件
        $(pwd)/api/v1alpha1/fuwu_types.go
		> 这个文件包含了对 Fuwu 这个 CRD 的定义，开发者可以在里面添加修改 Spec 和 Status
    -  CRD 控制器处理文件
        $(pwd)/controllers/fuwu_controller.go
        > 这个文件是控制器的处理逻辑，当集群中有 Fuwu 资源的变动（CRUD），都会触发 Reconcile 这个函数进行协调
    * **修改 spec**
    
        ```go
        // FuwuSpec defines the desired state of Fuwu
        type FuwuSpec struct {
            // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
            // Important: Run "make" to regenerate code after modifying this file
            
            // 以下是添加的字段信息
            Name        string `json:"name"`
            Description string `json:"description"`
        }
        ```
        
        ```bash
        # 修改 spec 参数后，执行以下命令，即可同步 crd spec yaml
        make && make install
        ```
        
        ```diff
        -	// Foo is an example field of Fuwu. Edit Fuwu_types.go to remove/update
        -	Foo string `json:"foo,omitempty"`
        +	Name        string `json:"name"`
        +	Description string `json:"description"`
        ```
        
        ```diff
                 spec:
                   description: FuwuSpec defines the desired state of Fuwu
                   properties:
        -            foo:
        -              description: Foo is an example field of Fuwu. Edit Fuwu_types.go to
        -                remove/update
        +            description:
                       type: string
        +            name:
        +              type: string
        +          required:
        +          - description
        +          - name
                   type: object
                 status:
                   description: FuwuStatus defines the observed state of Fuwu
        ```
    
    * 修改 Status
    
        ```go
        type FuwuStatus struct {
            // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
            // Important: Run "make" to regenerate code after modifying this file
        
            Status string `json:"status"`
        }
        ```
        
        碰到的问题
        ```bash
        ERRO[0000] fuwus.runner.basebit.me "fuwu-sample" not foundunable to update fuwu status  source="fuwu_controller.go:56"
        ```
        
        修复方法
        
        在 CRD 结构体上加上以下注释, 在CRD定义中启用状态子资源 
        
        // +kubebuilder:subresource:status
        
        ```bash
        # cat ./api/v1alpha1/fuwu_types.go
  
        // +kubebuilder:subresource:status
        // +kubebuilder:object:root=true
        ```
		
		```bash
		# 修改后，执行以下命令，即可同步 crd spec yaml
		make manifests
		```
		
		```diff
		diff --git a/config/crd/bases/runner.basebit.me_fuwus.yaml b/config/crd/bases/runner.basebit.me_fuwus.yaml
		index 66a42cc..c3aab3e 100644
		--- a/config/crd/bases/runner.basebit.me_fuwus.yaml
		+++ b/config/crd/bases/runner.basebit.me_fuwus.yaml
		@@ -15,6 +15,8 @@ spec:
			 plural: fuwus
			 singular: fuwu
		   scope: Namespaced
		+  subresources:
		+    status: {}
		   validation:
			 openAPIV3Schema:
			   description: Fuwu is the Schema for the fuwus API
		```
		
		

4. 安装 CRD
    ```bash
    # 安装 CRD
    make install
 
    # 查看创建的 CRD
    kubectl get crd apps.runner.basebit.me
    kubectl get crd fuwus.runner.basebit.me
    ```
    
5. 本地启动 controller
    
    ```bash
    # 启动 CRD controller
    make run
    ```
    
6. 部署 controller 到集群
    
    ```bash
    # 构建镜像
    make docker-build docker-push IMG=docker-reg.basebit.me:5000/base/runner-controller:v1alpha1
    
    # 部署到集群
    make deploy IMG=docker-reg.basebit.me:5000/base/runner-controller:v1alpha1
    ```
	
6. 创建一个 webhook
    
    > webhook server 需要 CA 证书
     
    如果需要在 FUWU CRUD 时进行操作合法性检查， 可以开发一个 webhook 实现。webhook 的脚手架一样可以用 kubebuilder 生成
    
    ```bash
    kubebuilder create webhook --group runner --version v1alpha1 --kind Fuwu --defaulting --programmatic-validation
    ```
	
7. 安装 Custom Resources 实例
	```bash
	kubectl apply -f config/samples/
	
	# 查看创建的实例
	➜ kubectl get fuwus.runner.basebit.me
	NAME          AGE
	fuwu-sample   3m2s
	
	➜ kubectl get apps.runner.basebit.me
	NAME         AGE
	app-sample   3m9s
	```

8. 卸载 CRDs
	```bash
	make uninstall
	```	

## controller 逻辑

1. controller 把轮询与事件监听都封装在这一个接口（Reconcile）里了，不需要关心怎么事件监听的.
2. 控制器的处理函数，每当集群中有 fuwu 资源的变动（CRUD），都会触发这个函数进行协调。

![image](https://user-images.githubusercontent.com/8086910/92990273-9451f300-f50d-11ea-9ba9-106e9087a01f.png)

### 如何同步自定义资源以及 K8s build-in 资源？
需要将自定义资源和想要 Watch 的 K8s build-in 资源的 GVKs 注册到 Scheme 上，Cache 会自动帮我们同步。

### Controller 的 Reconcile 方法是如何被触发的？
通过 Cache 里面的 Informer 获取资源的变更事件，然后通过两个内置的 Controller 以生产者消费者模式传递事件，最终触发 Reconcile 方法。

### Cache 的工作原理是什么？
GVK -> Informer 的映射，Informer 包含 Reflector 和 Indexer 来做事件监听和本地缓存。

### Implementing a controller

## References

- [The Kubebuilder Book](https://book.kubebuilder.io/introduction.html)
- [Status Subresource](https://book-v1.book.kubebuilder.io/basics/status_subresource.html)
- [进阶 K8s 高级玩家必备 | Kubebuilder：让编写 CRD 变得更简单](https://mp.weixin.qq.com/s/UzEcj2eXKM0m8f4XzZCYAA)